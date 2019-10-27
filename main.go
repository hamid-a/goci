package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"net/http"
	"os/exec"
	"time"
)

var e = echo.New()

func main() {
	e.HideBanner = true

	viper.SetConfigName("pipeline")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		e.Logger.Fatalf("Fatal error config file: %s", err.Error())
	}

	e.GET("/health", func(context echo.Context) error {
		return context.String(http.StatusOK, "It works!")
	})

	e.POST("/events", HandleEvents)

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.Fatal(e.Start(":" + viper.GetString("port")))
}

func HandleEvents(c echo.Context) error {

	token := c.Request().Header.Get("X-Gitlab-Token")
	service := c.QueryParam("service")
	accessToken := viper.GetString("token." + service)

	if accessToken != token {
		return c.String(http.StatusUnauthorized, "Unauthorized.")
	}

	go runPipeline(service)

	return c.String(http.StatusOK, "ok.")
}

func runPipeline(service string) {

	for _, pipeline := range viper.GetStringSlice(service + ".pipelines") {
		key := fmt.Sprintf("%s.%s.", service, pipeline)
		dir := viper.GetString(key + "directory")
		for _, command := range viper.GetStringSlice(key + "commands") {

			cmd := exec.Command("bash", "-c", command)
			cmd.Dir = dir

			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				e.Logger.Error(err)
			}
			fmt.Printf(
				"[%s] Service: %s, Command: %s, Path: %s\n%s",
				time.Now().Format("2006-01-02 15:04:05"),
				service,
				command,
				dir,
				out.String(),
			)
		}
	}
}

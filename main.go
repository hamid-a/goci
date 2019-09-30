package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func main()  {
	e := echo.New()

	e.GET("/hello", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
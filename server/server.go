package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

// New for instanctiate the server
func New() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":3000"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

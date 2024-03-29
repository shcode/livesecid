package server

import (
	"github.com/awcodify/livesecid/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New for instanctiate the server
func New() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/auth/google/sign_in", handler.GoogleSignIn)
	e.GET("/auth/google/callback", handler.GoogleCallbackAfterSuccess)

	e.Logger.Fatal(e.Start(":3000"))
}

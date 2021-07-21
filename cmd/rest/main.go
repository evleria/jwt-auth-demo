package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "gh/evleria/jwt-auth-demo/cmd/rest/docs"
)

// @title JWT Auth Demo Project
func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":1234"))
}

// Hello godoc
// @Summary Prints hello
// @Success 200 {string} Token "Hello"
// @Router / [get]
func hello(context echo.Context) error {
	return context.String(http.StatusOK, "Hello")
}
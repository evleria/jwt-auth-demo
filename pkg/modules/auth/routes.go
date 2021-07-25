package auth

import (
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddRoutes(group *echo.Group, database database.Database) {
	group.GET("/hello", hello)
}

// Hello godoc
// @Summary Prints hello
// @Success 200 {string} Token "Hello"
// @Router /auth/hello [get]
func hello(context echo.Context) error {
	return context.String(http.StatusOK, "Hello")
}

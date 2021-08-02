// Package middleware contains custom middlewares
package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/evleria/jwt-auth-demo/internal/service"
)

// Jwt produces echo.MiddlewareFunc with custom token validation
func Jwt(authService service.Auth) echo.MiddlewareFunc {
	return echoMiddleware.JWTWithConfig(echoMiddleware.JWTConfig{
		ParseTokenFunc: func(token string, ctx echo.Context) (interface{}, error) {
			return authService.ValidateAccessToken(ctx.Request().Context(), token)
		},
	})
}

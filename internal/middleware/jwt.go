package middleware

import (
	"github.com/evleria/jwt-auth-demo/internal/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Jwt(authService service.Auth) echo.MiddlewareFunc {
	return echoMiddleware.JWTWithConfig(echoMiddleware.JWTConfig{
		ParseTokenFunc: func(token string, ctx echo.Context) (interface{}, error) {
			return authService.ValidateAccessToken(ctx.Request().Context(), token)
		},
	})
}

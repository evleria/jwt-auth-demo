package auth

import (
	"github.com/evleria/jwt-auth-demo/internal/common/database"
	"github.com/evleria/jwt-auth-demo/internal/common/jwt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func AddModule(group *echo.Group, database database.Database, redis *redis.Client) {
	userRepository := NewUserRepository(database)
	tokenRepository := NewTokenRepository(redis)
	jwtMaker := jwt.NewJwtMaker(jwt.FromConfig())
	service := NewService(userRepository, tokenRepository, jwtMaker)
	controller := NewController(service)

	group.POST("/register", controller.Register)
	group.POST("/login", controller.Login)
	group.POST("/refresh", controller.Refresh)
}

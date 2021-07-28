package auth

import (
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/evleria/jwt-auth-demo/pkg/common/jwt"
	"github.com/evleria/jwt-auth-demo/pkg/common/kvstore"
	"github.com/labstack/echo/v4"
)

func AddModule(group *echo.Group, database database.Database, kvstore kvstore.KVStore) {
	userRepository := NewUserRepository(database)
	tokenRepository := NewTokenRepository(kvstore)
	jwtMaker := jwt.NewJwtMaker(jwt.FromConfig())
	service := NewService(userRepository, tokenRepository, jwtMaker)
	controller := NewController(service)

	group.POST("/register", controller.Register)
	group.POST("/login", controller.Login)
	group.POST("/refresh", controller.Refresh)
}

package auth

import (
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/evleria/jwt-auth-demo/pkg/common/jwt"
	"github.com/labstack/echo/v4"
)

func AddModule(group *echo.Group, database database.Database) {
	repository := NewRepository(database)
	jwtMaker := jwt.NewJwtMaker(jwt.FromConfig())
	service := NewService(repository, jwtMaker)
	controller := NewController(service)

	group.POST("/register", controller.Register)
	group.POST("/login", controller.Login)
}

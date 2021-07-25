package auth

import (
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/labstack/echo/v4"
)

func AddRoutes(group *echo.Group, database database.Database) {
	repository := NewRepository(database)
	service := NewService(repository)
	controller := NewController(service)

	group.POST("/register", controller.Register)
}

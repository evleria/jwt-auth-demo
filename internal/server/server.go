package server

import (
	"fmt"
	_ "github.com/evleria/jwt-auth-demo/docs"
	"github.com/evleria/jwt-auth-demo/internal/common/database"
	"github.com/evleria/jwt-auth-demo/internal/modules/auth"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server interface {
	Listen() error
}

type server struct {
	echo   *echo.Echo
	db     database.Database
	redis  *redis.Client
	config Config
}

func New(echo *echo.Echo, db database.Database, redis *redis.Client, config Config) Server {
	return &server{
		echo:   echo,
		db:     db,
		redis:  redis,
		config: config,
	}
}

func (s *server) Listen() error {
	s.initRoutes()
	return s.echo.Start(fmt.Sprintf(":%d", s.config.Port))
}

func (s *server) initRoutes() {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())

	auth.AddModule(
		s.echo.Group("/auth"),
		s.db,
		s.redis)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

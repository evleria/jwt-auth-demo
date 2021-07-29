package server

import (
	_ "github.com/evleria/jwt-auth-demo/docs"
	"github.com/evleria/jwt-auth-demo/internal/common/database"
	"github.com/evleria/jwt-auth-demo/internal/common/webserver"
	"github.com/evleria/jwt-auth-demo/internal/modules/auth"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server interface {
	Listen() error
}

type server struct {
	webserver webserver.WebServer
	db        database.Database
	redis     *redis.Client
	config    Config
}

func New(webServer webserver.WebServer, db database.Database, redis *redis.Client, config Config) Server {
	return &server{
		webserver: webServer,
		db:        db,
		redis:     redis,
		config:    config,
	}
}

func (s *server) Listen() error {
	s.initRoutes()
	return s.webserver.Listen(s.config.Port)
}

func (s *server) initRoutes() {
	engine := s.webserver.Engine()
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())

	auth.AddModule(
		engine.Group("/auth"),
		s.db,
		s.redis)

	engine.GET("/swagger/*", echoSwagger.WrapHandler)
}

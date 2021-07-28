package server

import (
	_ "github.com/evleria/jwt-auth-demo/cmd/rest/docs"
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/evleria/jwt-auth-demo/pkg/common/kvstore"
	"github.com/evleria/jwt-auth-demo/pkg/common/webserver"
	"github.com/evleria/jwt-auth-demo/pkg/modules/auth"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server interface {
	Listen() error
}

type server struct {
	webserver webserver.WebServer
	db        database.Database
	kvstore   kvstore.KVStore
	config    Config
}

func New(webServer webserver.WebServer, db database.Database, kvstore kvstore.KVStore, config Config) Server {
	return &server{
		webserver: webServer,
		db:        db,
		kvstore:   kvstore,
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
		s.kvstore)

	engine.GET("/swagger/*", echoSwagger.WrapHandler)
}

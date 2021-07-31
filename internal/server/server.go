package server

import (
	"fmt"
	_ "github.com/evleria/jwt-auth-demo/docs"
	"github.com/evleria/jwt-auth-demo/internal/controllers"
	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/repository"
	"github.com/evleria/jwt-auth-demo/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server interface {
	Listen(port int) error
}

type server struct {
	echo  *echo.Echo
	db    *pgx.Conn
	redis *redis.Client
}

func New(echo *echo.Echo, db *pgx.Conn, redis *redis.Client) Server {
	return &server{
		echo:  echo,
		db:    db,
		redis: redis,
	}
}

func (s *server) Listen(port int) error {
	s.initRoutes()
	return s.echo.Start(fmt.Sprintf(":%d", port))
}

func (s *server) initRoutes() {
	jwtMaker := jwt.NewJwtMakerFromConfig()
	userRepository := repository.NewUserRepository(s.db)
	tokenRepository := repository.NewTokenRepository(s.redis)
	authService := service.NewAuthService(userRepository, tokenRepository, jwtMaker)
	controller := controllers.NewAuthController(authService)

	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())

	authGroup := s.echo.Group("/auth")
	authGroup.POST("/register", controller.Register)
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/refresh", controller.Refresh)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/evleria/jwt-auth-demo/docs"
	"github.com/evleria/jwt-auth-demo/internal/config"
	"github.com/evleria/jwt-auth-demo/internal/handler"
	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/repository"
	"github.com/evleria/jwt-auth-demo/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"time"
)

var cfgPath string

func initFlags() {
	flag.StringVar(&cfgPath, "config", "./configs/.env", "The application configuration")
	flag.Parse()
}

// @title JWT Auth Demo Project
func main() {
	initFlags()

	err := config.Load(cfgPath)
	if err != nil {
		log.Println(err)
	}

	dbUrl := getPostgresConnectionString()
	db, err := pgx.Connect(context.Background(), dbUrl)
	check(err)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getRedisAddress(),
		Password: config.GetString("REDIS_PASSWORD", ""),
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	check(err)

	e := echo.New()
	initRoutes(e, db, redisClient)

	port := fmt.Sprintf(":%d", config.GetInt("PORT", 5000))
	check(e.Start(port))
}

func initRoutes(e *echo.Echo, db *pgx.Conn, redisClient *redis.Client) {
	jwtMaker := jwt.NewJwtMakerFromConfig()
	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokenRepository(redisClient)
	authService := service.NewAuthService(userRepository, tokenRepository, jwtMaker)
	authHandler := handler.NewAuthHandler(authService)

	e.Use(middleware.Logger())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: time.Second * 5}))
	e.Use(middleware.Recover())

	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func getPostgresConnectionString() string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.GetString("POSTGRES_USER", "postgres"),
		config.GetString("POSTGRES_PASSWORD", ""),
		config.GetString("POSTGRES_HOST", "127.0.0.1"),
		config.GetInt("POSTGRES_PORT", 5432),
		config.GetString("POSTGRES_DB", "postgres"),
	)

	if config.GetBool("POSTGRES_SSL_DISABLE", false) {
		conn += "?sslmode=disable"
	}

	return conn
}

func getRedisAddress() string {
	return fmt.Sprintf("%s:%d",
		config.GetString("REDIS_HOST", "localhost"),
		config.GetInt("REDIS_PORT", 6379),
	)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

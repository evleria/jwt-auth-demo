package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/evleria/jwt-auth-demo/internal/common/config"
	"github.com/evleria/jwt-auth-demo/internal/common/database"
	"github.com/evleria/jwt-auth-demo/internal/common/kvstore"
	"github.com/evleria/jwt-auth-demo/internal/common/webserver"
	"github.com/evleria/jwt-auth-demo/internal/server"
	"github.com/go-redis/redis/v8"
	"log"
)

var cfgPath string

func initFlags() {
	flag.StringVar(&cfgPath, "config", "./configs/.env", "The application configuration")
	flag.Parse()
}

// @title JWT Auth Demo Project
func main() {
	initFlags()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := config.Load(cfgPath)
	if err != nil {
		log.Println(err)
	}

	dbUrl := getPostgresConnectionString()

	db, err := database.New(ctx, dbUrl)
	check(err)

	kvstore, err := kvstore.New(ctx, &redis.Options{
		Addr:     getRedisAddress(),
		Password: config.GetString("REDIS_PASSWORD", ""),
	})
	check(err)

	webserver := webserver.New()

	server := server.New(webserver, db, kvstore, server.Config{
		Port: config.GetInt("PORT", 5000),
	})

	check(server.Listen())
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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/evleria/jwt-auth-demo/pkg/common/config"
	"github.com/evleria/jwt-auth-demo/pkg/common/database"
	"github.com/evleria/jwt-auth-demo/pkg/common/webserver"
	"github.com/evleria/jwt-auth-demo/pkg/server"
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

	dbUrl := GetPostgresConnectionString()

	db, err := database.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	webserver := webserver.New()

	server := server.New(webserver, db, server.Config{
		Port: config.GetInt("PORT", 5000),
	})

	log.Fatal(server.Listen())
}

func GetPostgresConnectionString() string {
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

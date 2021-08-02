// Package config encapsulates work with environment variables
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Load loads .env file by given path
func Load(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	if err := godotenv.Load(absPath); err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	return nil
}

// GetString provides string config variable by key and given fallback
func GetString(key, fallback string) string {
	return lookup(key, fallback)
}

// GetInt provides int config variable by key and given fallback
func GetInt(key string, fallback int) int {
	value := lookup(key, "")
	if value, err := strconv.Atoi(value); err == nil {
		return value
	}
	return fallback
}

// GetBool provides bool config variable by key and given fallback
func GetBool(key string, fallback bool) bool {
	value := lookup(key, "")
	if value, err := strconv.ParseBool(value); err == nil {
		return value
	}
	return fallback
}

// GetDuration provides time.Duration config variable by key and given fallback
func GetDuration(key string, fallback time.Duration) time.Duration {
	value := lookup(key, "")
	if value, err := time.ParseDuration(value); err == nil {
		return value
	}
	return fallback
}

func lookup(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

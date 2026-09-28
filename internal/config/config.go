package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Database_url string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		panic("DATABASE_URL is required")
	}

	return Config{
		Port:         port,
		Database_url: database_url,
	}
}

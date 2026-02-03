package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv  string
	AppPort string
	DSN     string
}

func SetupEnv() (AppConfig, error) {

	// Load .env only in development
	appEnv := os.Getenv("APP_ENV")
	if  appEnv != "production" {
		if err := godotenv.Load(); err != nil {
			return AppConfig{}, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, fmt.Errorf("HTTP_PORT environment variable is not set")
	}

	Dsn := os.Getenv("DSN")

	if len(Dsn)< 1 {
		return  AppConfig{}, errors.New("Dsn variable not found")
	}

	return AppConfig{
		AppEnv: appEnv,
		AppPort: httpPort,
		DSN: Dsn,
	}, nil

}

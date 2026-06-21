package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	AppEnv       string
	JWTSecret    string
	DatabasePath string
	LogLevel     string
}

func NewConfig() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es requerido")
	}
	if len(jwtSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		return Config{}, fmt.Errorf("DATABASE_PATH es requerido")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		Port:         port,
		AppEnv:       appEnv,
		JWTSecret:    jwtSecret,
		DatabasePath: databasePath,
		LogLevel:     logLevel,
	}, nil
}

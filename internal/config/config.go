package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName            string
	AppEnv             string
	AppPort            string
	DBConnection       string
	DBHost             string
	DBPort             string
	DBDatabase         string
	DBUsername         string
	DBPassword         string
	JWTSecret          string
	JWTExpirationHours int
}

func LoadConfig() *Config {
	// Try loading from .env if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or unable to load, reading from environment")
	}

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		jwtExpHours = 72
	}

	return &Config{
		AppName:            getEnv("APP_NAME", "EcoPoints-Go-API"),
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBConnection:       getEnv("DB_CONNECTION", "mysql"),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBDatabase:         getEnv("DB_DATABASE", "ecopoints"),
		DBUsername:         getEnv("DB_USERNAME", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "root"),
		JWTSecret:          getEnv("JWT_SECRET", "ecopoints_jwt_secret_key_2026_super_secure_key"),
		JWTExpirationHours: jwtExpHours,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

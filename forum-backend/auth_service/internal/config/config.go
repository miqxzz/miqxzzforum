package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	MigrationsPath string
	JWTSecret      string
}

func LoadConfig() (Config, error) {

	err := godotenv.Load()
	if err != nil {

	}

	cfg := Config{
		Port:           getEnv("AUTH_SERVICE_PORT", ":8080"),
		DBPath:         getEnv("DB_PATH", "../db/forum.db"),
		MigrationsPath: getEnv("AUTH_SERVICE_MIGRATIONS_PATH", "C:\\forum-project\\forum-backend\\auth_service\\migrations"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
	}
	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

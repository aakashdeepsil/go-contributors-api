package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	MongoDB struct {
		URI      string
		Database string
	}
	Redis struct {
		URL      string
		Password string
		DB       int
	}
	RateLimit struct {
		PerSecond int
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
	}

	// MongoDB configuration
	config.MongoDB.URI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	config.MongoDB.Database = getEnv("MONGODB_DATABASE", "contributors_db")

	// Redis configuration
	config.Redis.URL = getEnv("REDIS_URL", "localhost:6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

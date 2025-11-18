package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost   string
	ServerPort   string
	DataFile     string
	CheckTimeout time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("not found .env file: %s", err)
	}

	cfg := &Config{
		ServerHost:   getEnv("SERVER_HOST", "localhost"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DataFile:     getEnv("DATA_FILE", "data.json"),
		CheckTimeout: parseDuration(getEnv("CHECK_TIMEOUT", "10s")),
	}

	return cfg, nil
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}

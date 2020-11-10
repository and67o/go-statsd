package config

import (
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("logger/.env"); err != nil {
		panic("No .env file found")
	}
}

func New() *Config {
	return &Config{
		Stats: StatsDConfig{
			Host:  getEnv("STATSD_HOST", "127.0.0.1"),
			Port: getEnvAsInt("STATSD_PORT", 8125),
		},
	}
}

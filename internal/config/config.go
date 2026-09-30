package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultConnstr = "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	defaultAddress = ":8080"
)

type Config struct {
	Connstr string
	Address string
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load(files ...string) (*Config, error) {
	godotenv.Load(files...)
	return &Config{
			Connstr: envOrDefault("CONNSTR", defaultConnstr),
			Address: envOrDefault("ADDRESS", defaultAddress),
		},
		nil
}

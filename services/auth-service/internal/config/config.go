package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	NatsURL     string
	JwtSecret   string
}

func Load() *Config {
	db := os.Getenv("DATABASE_URL")
	nats := os.Getenv("NATS_URL")
	jwt := os.Getenv("JWT_SECRET")

	if db == "" || nats == "" || jwt == "" {
		log.Fatal("Missing required environment variables")
	}

	return &Config{
		DatabaseURL: db,
		NatsURL:     nats,
		JwtSecret:   jwt,
	}
}

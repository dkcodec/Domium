package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	NatsURL     string
	JwtSecret   string
	GRPCPort    int
}

func Load() *Config {
	db := os.Getenv("DATABASE_URL")
	nats := os.Getenv("NATS_URL")
	jwt := os.Getenv("JWT_SECRET")
	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal("Invalid GRPC_PORT value")
	}

	if db == "" || nats == "" || jwt == "" || grpcPort == 0 {
		log.Fatal("Missing required environment variables")
	}

	return &Config{
		DatabaseURL: db,
		NatsURL:     nats,
		JwtSecret:   jwt,
		GRPCPort:    grpcPort,
	}
}

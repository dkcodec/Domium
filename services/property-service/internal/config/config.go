package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	NatsURL     string
	GRPCPort    int
}

func Load() *Config {
	db := os.Getenv("DATABASE_URL")
	nats := os.Getenv("NATS_URL")
	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal("Invalid GRPC_PORT value")
	}

	if db == "" || nats == "" {
		log.Fatal("Missing required environment variables")
	}

	return &Config{
		DatabaseURL: db,
		NatsURL:     nats,
		GRPCPort:    grpcPort,
	}
}

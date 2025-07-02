package main

import (
	"fmt"
	"log"
	"net"
	"property-service/internal/config"
	"property-service/internal/handler"
	"property-service/internal/proto"
	"property-service/internal/repository"
	"property-service/internal/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load(".env")

	cfg := config.Load()

	db := repository.NewPostgresDB(cfg.DatabaseURL)
	defer db.Close()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	repo := repository.NewPostgresRepo(db)
	service := service.New(repo, nc)
	handler := handler.New(service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPropertyServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
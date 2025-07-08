package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"

	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load(".env") // Загружаем переменные окружения из файла

	// Загружаем конфиг из .env или переменных окружения
	cfg := config.Load()

	// Подключаемся к PostgreSQL через репозиторий
	db := postgres.NewPostgres(cfg.DatabaseURL)	
	defer db.Close()

	// Подключаемся к NATS для публикации событий
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Инициализируем usecase слой (бизнес-логику)
	userRepo := postgres.NewUserRepo(db)
	authService := service.NewAuthService(userRepo, nc, cfg.JwtSecret)

	// gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем gRPC хендлер, передаём бизнес-логику
	handler.RegisterAuthServer(grpcServer, authService)

	// Слушаем порт
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
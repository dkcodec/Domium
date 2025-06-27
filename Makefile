.PHONY: up down build proto migrate

# Запуск всех сервисов
up:
	docker-compose up --build

# Остановка всех сервисов
down:
	docker-compose down

# Сборка proto-файлов
proto:
	protoc --go_out=. --go-grpc_out=. services/auth-service/proto/user.proto

# Миграция users таблицы
migrate:
	docker exec -i rent-postgres-1 psql -U postgres -d auth_db < services/auth-service/migrations/001_init.sql

# Пересборка auth-service
build-auth:
	docker-compose build auth-service

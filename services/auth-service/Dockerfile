# Dockerfile
FROM golang:1.23.10

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

EXPOSE 50051

CMD ["./auth-service"]

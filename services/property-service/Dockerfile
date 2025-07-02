FROM golang:1.23.10

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o property ./cmd/main.go

CMD ["./property"]

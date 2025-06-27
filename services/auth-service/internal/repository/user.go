package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID          string
	PhoneNumber string
	PasswordHash string
	FullName    string
	CreatedAt   time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByPhone(ctx context.Context, phone string) (*User, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgres(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, phone_number, password_hash, full_name, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.PhoneNumber, user.PasswordHash, user.FullName, user.CreatedAt,
	)
	return err
}

func (r *PostgresRepo) FindByPhone(ctx context.Context, phone string) (*User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, phone_number, password_hash, full_name, created_at
		 FROM users WHERE phone_number = $1`, phone)

	var user User
	err := row.Scan(&user.ID, &user.PhoneNumber, &user.PasswordHash, &user.FullName, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

package postgres

import (
	"auth-service/internal/repository"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type AuthPostgresRepo struct {
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

func NewUserRepo(db *sql.DB) *AuthPostgresRepo {
	return &AuthPostgresRepo{db: db}
}

func (r *AuthPostgresRepo) CreateUser(ctx context.Context, user *repository.User) (string, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, is_email_verified, verification_code, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.Email, user.PasswordHash, user.IsEmailVerified, user.VerificationCode, user.CreatedAt, user.UpdatedAt,
	)
	return user.ID, err
}

func (r *AuthPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	row := r.db.QueryRowContext(ctx,
	`SELECT id, email, is_email_verified, password_hash, verification_code, created_at FROM users WHERE email = $1`, email)
	
	var user repository.User
	err := row.Scan(&user.ID, &user.Email, &user.IsEmailVerified, &user.PasswordHash, &user.VerificationCode, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *AuthPostgresRepo) GetUserByID(ctx context.Context, id string) (*repository.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, is_email_verified, verification_code, created_at FROM users WHERE id = $1`, id)

	var user repository.User
	err := row.Scan(&user.ID, &user.Email, &user.IsEmailVerified, &user.VerificationCode, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *AuthPostgresRepo) UpdateUser(ctx context.Context, user *repository.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email = $1, is_email_verified = $2, verification_code = $3, updated_at = $4 WHERE id = $5`,
		user.Email, user.IsEmailVerified, user.VerificationCode, user.UpdatedAt, user.ID)
	return err
}

func (r *AuthPostgresRepo) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}


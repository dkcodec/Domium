package repository

import "time"

type User struct {
	ID               string
	Email            string
	PasswordHash     string
	IsEmailVerified  bool
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UpdateUserInput struct {
	ID               string
	Email            string
	IsEmailVerified  bool
	VerificationCode string
	UpdatedAt        time.Time
}
package repository

import "context"

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	VerifyEmailCode(ctx context.Context, email, code string) error
	SetEmailVerified(ctx context.Context, email string) error
}
package service

import (
	"auth-service/internal/repository"
	"auth-service/internal/repository/postgres"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *postgres.AuthPostgresRepo
	nats      *nats.Conn
	secretKey string
}

func NewAuthService(repo *postgres.AuthPostgresRepo, nats *nats.Conn, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		nats:      nats,
		secretKey: secret,
	}
}

func (s *AuthService) ResendVerificationCode(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	user.VerificationCode = uuid.New().String()
	user.UpdatedAt = time.Now()
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return s.nats.Publish("user.resend_verification", []byte(user.Email))
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}

	verificationCode := uuid.New().String()
	user := &repository.User{
		Email:            email,
		PasswordHash:     string(hash),
		IsEmailVerified:  false,
		VerificationCode: verificationCode,
	}
	if _, err := s.repo.CreateUser(ctx, user); err != nil {
		return "", "", "", err
	}
	access, err := s.generateJWT(user.ID, "access", 15*time.Minute)
	if err != nil {
		return "", "", "", err
	}
	refresh, err := s.generateJWT(user.ID, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", "", err
	}

	_ = s.nats.Publish("user.registered", []byte(user.ID))
	return access, refresh, user.ID, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	fmt.Println("user", user)
	fmt.Println("err", err)
	if err != nil || user == nil {
		return "", "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", "", errors.New("invalid credentials")
	}
	access, err := s.generateJWT(user.ID, "access", 15*time.Minute)
	if err != nil {
		return "", "", "", err
	}
	refresh, err := s.generateJWT(user.ID, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", "", err
	}
	return access, refresh, user.ID, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.VerificationCode != code {
		return errors.New("invalid verification code")
	}
	user.IsEmailVerified = true
	user.UpdatedAt = time.Now()
	err = s.repo.UpdateUser(ctx, user)
	return err
}

func (s *AuthService) generateJWT(userID, tokenType string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiresIn).Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) GetUser(ctx context.Context, id string) (*repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *AuthService) UpdateUser(ctx context.Context, user *repository.User) error {
	err := s.repo.UpdateUser(ctx, user)
	return err
}

func (s *AuthService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

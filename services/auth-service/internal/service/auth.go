package service

import (
	"auth-service/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	nats      *nats.Conn
	secretKey string
}

func NewAuthService(repo repository.UserRepository, nats *nats.Conn, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		nats:      nats,
		secretKey: secret,
	}
}

func (s *AuthService) Register(ctx context.Context, phone, password, fullName string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	user := &repository.User{
		PhoneNumber:  phone,
		PasswordHash: string(hash),
		FullName:     fullName,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", "", err
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", "", err
	}

	_ = s.nats.Publish("user.registered", []byte(user.ID))
	return token, user.ID, nil
}

func (s *AuthService) Login(ctx context.Context, phone, password string) (string, string, error) {
	user, err := s.repo.FindByPhone(ctx, phone)
	if err != nil || user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", "", err
	}
	return token, user.ID, nil
}

func (s *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

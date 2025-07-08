package handler

import (
	"auth-service/internal/proto"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"context"
	"time"

	"google.golang.org/grpc"
)

type GRPCHandler struct {
	proto.UnimplementedAuthServiceServer
	auth *service.AuthService
}

func RegisterAuthServer(s *grpc.Server, auth *service.AuthService) {
	proto.RegisterAuthServiceServer(s, &GRPCHandler{auth: auth})
}

func (h *GRPCHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	access, refresh, userID, err := h.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{
		Id:            userID,
		AccessToken:   access,
		RefreshToken:  refresh,
	}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	access, refresh, _, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil	
}

func (h *GRPCHandler) VerifyEmail(ctx context.Context, req *proto.VerifyEmailRequest) (*proto.VerifyEmailResponse, error) {
	err := h.auth.VerifyEmail(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &proto.VerifyEmailResponse{}, nil
}

func (h *GRPCHandler) ResendVerificationCode(ctx context.Context, req *proto.ResendVerificationCodeRequest) (*proto.ResendVerificationCodeResponse, error) {
	err := h.auth.ResendVerificationCode(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &proto.ResendVerificationCodeResponse{}, nil
}

func (h *GRPCHandler) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := h.auth.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{
		User: &proto.User{
			Id:              user.ID,
			Email:           user.Email,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *GRPCHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	user := &repository.User{
		ID:               req.Id,
		Email:            req.Email,
		IsEmailVerified:  req.IsEmailVerified,
		VerificationCode: req.VerificationCode,
		UpdatedAt:        time.Now(),
	}
	return &proto.UpdateUserResponse{}, h.auth.UpdateUser(ctx, user)
}

func (h *GRPCHandler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return &proto.DeleteUserResponse{}, h.auth.DeleteUser(ctx, req.Id)
}

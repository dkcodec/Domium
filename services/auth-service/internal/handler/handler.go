package handler

import (
	"auth-service/internal/proto"
	"auth-service/internal/service"
	"context"

	"google.golang.org/grpc"
)

type GRPCHandler struct {
	proto.UnimplementedAuthServer
	auth *service.AuthService
}

func RegisterAuthServer(s *grpc.Server, auth *service.AuthService) {
	proto.RegisterAuthServer(s, &GRPCHandler{auth: auth})
}

func (h *GRPCHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {
	token, userID, err := h.auth.Register(ctx, req.PhoneNumber, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{Token: token, UserId: userID}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	token, userID, err := h.auth.Login(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{Token: token, UserId: userID}, nil
}

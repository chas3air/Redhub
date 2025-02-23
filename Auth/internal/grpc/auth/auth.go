package grpcauth

import (
	"auth/internal/domain/interfaces"
	authservice "auth/internal/services/auth"
	"auth/internal/storage"
	"context"
	"errors"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth interfaces.Auth
}

func Register(gRPC *grpc.Server, auth interfaces.Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	email := in.GetEmail()
	password := in.GetPassword()
	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	s_aid := in.GetAppId()
	if s_aid == "" {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	aid_uuid, err := uuid.Parse(s_aid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "app_id must be uuid")
	}

	token, err := s.auth.Login(ctx, email, password, aid_uuid)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	email := in.GetEmail()
	password := in.GetPassword()
	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	uid, err := s.auth.Register(ctx, email, password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &authv1.RegisterResponse{
		UserId: uid.String(),
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	uuid_uid, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "user_id must be uuid")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, uuid_uid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

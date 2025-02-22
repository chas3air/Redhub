package grpcauth

import (
	"auth/internal/domain/interfaces"
	"context"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth interfaces.Auth
}

func Register(gRPC *grpc.Server, auth interfaces.Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("unimplemented")
}

func (s *serverAPI) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("unimplemented")
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("unimplemented")
}

package interfaces

import (
	"context"

	authv1 "github.com/chas3air/protos/gen/go/github.chas3air.protos.auth"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID string) (token string, err error)
	Register(ctx context.Context, email string, password string) (user_id string, err error)
	IsAdmin(ctx context.Context, user_id string) (isAdmin bool, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func New(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

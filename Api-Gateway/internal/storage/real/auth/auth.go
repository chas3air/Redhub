package authstorage

import (
	"apigateway/internal/domain/models"
	authprofiles "apigateway/internal/domain/profiles/auth_profiles"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthStorage struct {
	log        *slog.Logger
	ServerHost string
	ServerPort int
}

func New(log *slog.Logger, host string, port int) *AuthStorage {
	return &AuthStorage{
		log:        log,
		ServerHost: host,
		ServerPort: port,
	}
}

func (as *AuthStorage) Login(ctx context.Context, email string, password string, app_id uuid.UUID) (token string, err error) {
	const op = "service.auth.login"
	log := as.log.With(
		slog.String("op", op),
	)
	_ = log

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", as.ServerHost, as.ServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := authv1.NewAuthClient(conn)
	res, err := c.Login(ctx,
		&authv1.LoginRequest{Email: email,
			Password: password,
			AppId:    app_id.String(),
		},
	)
	if err != nil {
		log.Warn("failed to get users", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetToken(), nil
}

func (as *AuthStorage) Register(ctx context.Context, user models.User) (err error) {
	const op = "service.auth.register"
	log := as.log.With(
		slog.String("op", op),
	)
	_ = log

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", as.ServerHost, as.ServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := authv1.NewAuthClient(conn)

	proto_user, _ := authprofiles.UsrToProtoUsr(user)

	_, err = c.Register(ctx, &authv1.RegisterRequest{
		User: proto_user,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (as *AuthStorage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "service.auth.isAdmin"
	log := as.log.With(
		slog.String("op", op),
	)
	_ = log

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", as.ServerHost, as.ServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := authv1.NewAuthClient(conn)
	res, err := c.IsAdmin(ctx, &authv1.IsAdminRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res.IsAdmin, nil
}

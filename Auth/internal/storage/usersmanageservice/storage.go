package usersmanageservice

import (
	"auth/internal/domain/models"
	"auth/internal/domain/profiles"
	"context"
	"fmt"
	"log/slog"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UsersManageService struct {
	log         *slog.Logger
	ServiceHost string
	ServicePort int
}

func New(log *slog.Logger, serviceHost string, servicePort int) *UsersManageService {
	return &UsersManageService{
		log:         log,
		ServiceHost: serviceHost,
		ServicePort: servicePort,
	}
}

// GetUsers implements interfaces.UsersStorage.
func (u *UsersManageService) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "usersmanageservice.getUsers"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	res, err := c.GetUsers(ctx, nil)
	if err != nil {
		log.Warn("failed to get users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var resUsers = make([]models.User, 0, len(res.GetUsers()))
	for _, pbUser := range res.GetUsers() {
		user, err := profiles.ProtoUsrToUsr(pbUser)
		if err != nil {
			log.Warn("failed to convert proto user to model user", slog.String("error", err.Error()))
			continue
		}
		resUsers = append(resUsers, user)
	}

	return resUsers, nil
}

// GetUserById implements interfaces.UsersStorage.
func (u *UsersManageService) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "usersmanageservice.getUserById"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	res, err := c.GetUserById(ctx, &umv1.GetUserByIdRequest{Id: uid.String()})
	if err != nil {
		log.Warn("failed to get user by ID", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	resUser, err := profiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Warn("failed to convert proto user to model user", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return resUser, nil
}

// GetUserByEmail implements interfaces.UsersStorage.
func (u *UsersManageService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "usersmanageservice.getUserByEmail"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	res, err := c.GetUserByEmail(ctx, &umv1.GetUserByEmailRequest{Email: email})
	if err != nil {
		log.Warn("failed to get user by email", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	resUser, err := profiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Warn("failed to convert proto user to model user", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return resUser, nil
}

// Insert implements interfaces.UsersStorage.
func (u *UsersManageService) Insert(ctx context.Context, user models.User) error {
	const op = "usersmanageservice.insert"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	userForInsert, err := profiles.UsrToProroUsr(user)
	if err != nil {
		log.Warn("failed to convert model user to proto user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = c.Insert(ctx, &umv1.InsertRequest{
		User: userForInsert,
	})
	if err != nil {
		log.Warn("failed to insert user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Update implements interfaces.UsersStorage.
func (u *UsersManageService) Update(ctx context.Context, uid uuid.UUID, user models.User) error {
	const op = "usersmanageservice.update"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	userForUpdate, err := profiles.UsrToProroUsr(user)
	if err != nil {
		log.Warn("failed to convert model user to proto user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = c.Update(ctx, &umv1.UpdateRequest{
		Id:   uid.String(),
		User: userForUpdate,
	})
	if err != nil {
		log.Warn("failed to update user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete implements interfaces.UsersStorage.
func (u *UsersManageService) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "usersmanageservice.delete"
	log := u.log.With(slog.String("op", op))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", u.ServiceHost, u.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := umv1.NewUsersManagerClient(conn)
	res, err := c.Delete(ctx, &umv1.DeleteRequest{
		Id: uid.String(),
	})
	if err != nil {
		log.Warn("failed to delete user", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	resUser, err := profiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Warn("failed to convert proto user to model user", slog.String("error", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return resUser, nil
}

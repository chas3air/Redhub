package usersmanager

import (
	"context"
	"errors"
	"log/slog"
	"usersManageService/internal/domain/interfaces/usersservice"
	"usersManageService/internal/domain/profiles"
	"usersManageService/internal/services"
	"usersManageService/pkg/lib/logger/sl"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type serverAPI struct {
	umv1.UnimplementedUsersManagerServer
	userManager usersservice.UsersManager
	log         *slog.Logger
}

func Register(grpc *grpc.Server, userManager usersservice.UsersManager, log *slog.Logger) {
	umv1.RegisterUsersManagerServer(grpc, &serverAPI{userManager: userManager, log: log})
}

func (s *serverAPI) GetUsers(ctx context.Context, req *umv1.GetUsersRequest) (*umv1.GetUsersResponse, error) {
	const op = "grpc.users.getUsers"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	app_users, err := s.userManager.GetUsers(ctx)
	if err != nil {
		log.Error("Error retrieving users", sl.Err(err))
		return nil, status.Error(codes.Internal, "error retrieving users")
	}

	resp_users := make([]*umv1.User, 0, len(app_users))
	for _, user := range app_users {
		profiles_user, err := profiles.UsrToProtoUsr(user)
		if err != nil {
			log.Error("Wrong structure", sl.Err(err))
			continue
		}

		resp_users = append(resp_users, profiles_user)
	}

	return &umv1.GetUsersResponse{
		Users: resp_users,
	}, nil
}

func (s *serverAPI) GetUserById(ctx context.Context, req *umv1.GetUserByIdRequest) (*umv1.GetUserByIdResponse, error) {
	const op = "grpc.users.getUserById"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	id_s := req.GetId()
	if id_s == "" {
		log.Error("Failed to get id", sl.Err(errors.New("required parametr id")))
		return nil, status.Error(codes.InvalidArgument, "required parametr id")
	}

	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid")
		return nil, status.Error(codes.InvalidArgument, "Invalid id, must be uuid")
	}

	requested_user, err := s.userManager.GetUserById(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("User with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "user with current id not found")
		}

		log.Error("failed to retrieve user by id", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve user by id")
	}

	profiled_user, err := profiles.UsrToProtoUsr(requested_user)
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "wrong structure")
	}

	return &umv1.GetUserByIdResponse{
		User: profiled_user,
	}, nil
}

func (s *serverAPI) GetUserByEmail(ctx context.Context, req *umv1.GetUserByEmailRequest) (*umv1.GetUserByEmailResponse, error) {
	const op = "grpc.users.getUserByEmail"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	req_email := req.GetEmail()
	if req_email == "" {
		log.Error("Failed to get email", sl.Err(errors.New("failed to get email")))
		return nil, status.Error(codes.InvalidArgument, "failed to get email")
	}

	requested_user, err := s.userManager.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("User with current email not found", sl.Err(err))
			return nil, status.Error(codes.InvalidArgument, "user with current email not found")
		}

		log.Error("Failed to retrieve user by email", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve user by email")
	}

	profiled_user, err := profiles.UsrToProtoUsr(requested_user)
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "wrong structure")
	}

	return &umv1.GetUserByEmailResponse{
		User: profiled_user,
	}, nil
}

func (s *serverAPI) Insert(ctx context.Context, req *umv1.InsertRequest) (*umv1.InsertResponse, error) {
	const op = "grpc.users.insertArticle"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	if req.GetUser() == nil {
		log.Error("User is required", sl.Err(errors.New("user is required")))
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	parsedUser, err := profiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "wrong structure")
	}

	err = s.userManager.Insert(ctx, parsedUser)
	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(err))
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Error("Failed to insert user", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "failed to insert user")
	}

	return nil, nil
}

func (s *serverAPI) Update(ctx context.Context, req *umv1.UpdateRequest) (*umv1.UpdateResponse, error) {
	const op = "grpc.users.updateArticle"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	if req.GetUser() == nil {
		log.Error("User is required", sl.Err(errors.New("user is required")))
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	parsedUser, err := profiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "wrong structure")
	}

	if req.GetId() == "" {
		log.Error("Id is required", sl.Err(errors.New("id is required")))
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	parsedUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("Invalid id, must be uuid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "Invalid id, must be uuid")
	}

	err = s.userManager.Update(ctx, parsedUUID, parsedUser)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("User with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "user with current id not found")
		}

		log.Error("Failed to update user", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return nil, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *umv1.DeleteRequest) (*umv1.DeleteResponse, error) {
	const op = "grpc.users.deleteArticle"
	log := s.log.With(
		slog.String("op", op),
	)
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	if req.GetId() == "" {
		log.Error("Id is required", sl.Err(errors.New("id is required")))
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	parsedUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("Invalid id, must be uuid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id, must be uuid")
	}

	user, err := s.userManager.Delete(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("User with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "user with current id not found")
		}

		log.Error("Failed to delete user", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	profiled_user, err := profiles.UsrToProtoUsr(user)
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "wrong structure")
	}

	return &umv1.DeleteResponse{
		User: profiled_user,
	}, nil
}

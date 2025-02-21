package usersmanager

import (
	"context"
	"log"
	"usersManageService/internal/domain/interfaces/usersservice"
	"usersManageService/internal/domain/profiles"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type serverAPI struct {
	umv1.UnimplementedUsersManagerServer
	userManager usersservice.UsersManager
}

func Register(grpc *grpc.Server, userManager usersservice.UsersManager) {
	umv1.RegisterUsersManagerServer(grpc, &serverAPI{userManager: userManager})
}

func (s *serverAPI) GetUsers(ctx context.Context, req *umv1.GetUsersRequest) (*umv1.GetUsersResponse, error) {
	app_users, err := s.userManager.GetUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve app users")
	}

	resp_users := make([]*umv1.User, 0, len(app_users))
	for _, user := range app_users {
		profiles_user, err := profiles.UsrToProroUsr(user)
		if err != nil {
			log.Println("error user:", user)
		}
		resp_users = append(resp_users, profiles_user)
	}

	return &umv1.GetUsersResponse{
		Users: resp_users,
	}, nil
}

func (s *serverAPI) GetUserById(ctx context.Context, req *umv1.GetUserByIdRequest) (*umv1.GetUserByIdResponse, error) {
	parsedUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	requested_user, err := s.userManager.GetUserById(ctx, parsedUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user no found")
	}
	profiled_user, err := profiles.UsrToProroUsr(requested_user)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}

	return &umv1.GetUserByIdResponse{
		User: profiled_user,
	}, nil
}

func (s *serverAPI) GetUserByEmail(ctx context.Context, req *umv1.GetUserByEmailRequest) (*umv1.GetUserByEmailResponse, error) {
	requested_user, err := s.userManager.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	profiled_user, err := profiles.UsrToProroUsr(requested_user)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}

	return &umv1.GetUserByEmailResponse{
		User: profiled_user,
	}, nil
}

func (s *serverAPI) Insert(ctx context.Context, req *umv1.InsertRequest) (*umv1.InsertResponse, error) {
	parsedUser, err := profiles.ProtoUsrToUsr((req.GetUser()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	err = s.userManager.Insert(ctx, parsedUser)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	return nil, nil
}

func (s *serverAPI) Update(ctx context.Context, req *umv1.UpdateRequest) (*umv1.UpdateResponse, error) {
	parsedUser, err := profiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	parsedUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	err = s.userManager.Update(ctx, parsedUUID, parsedUser)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	return nil, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *umv1.DeleteRequest) (*umv1.DeleteResponse, error) {
	parsedUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	user, err := s.userManager.Delete(ctx, parsedUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	profiled_user, err := profiles.UsrToProroUsr(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}

	return &umv1.DeleteResponse{
		User: profiled_user,
	}, nil
}

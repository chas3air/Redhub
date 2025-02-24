package authprofiles

import (
	"apigateway/internal/domain/models"
	"time"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UsrToProtoUsr(user models.User) (*authv1.User, error) {
	var birthday *timestamppb.Timestamp
	if !user.Birthday.IsZero() {
		birthday = timestamppb.New(user.Birthday)
	}

	return &authv1.User{
		Id:          user.Id.String(),
		Email:       user.Email,
		Password:    user.Password,
		Role:        user.Role,
		Nick:        user.Nick,
		Description: user.Description,
		Birthday:    birthday,
	}, nil
}

func ProtoUsrToUsr(proto_usr *authv1.User) (models.User, error) {
	parsedUUID, err := uuid.Parse(proto_usr.GetId())
	if err != nil {
		return models.User{}, err
	}

	var birthday time.Time
	if proto_usr.GetBirthday() != nil {
		birthday = proto_usr.GetBirthday().AsTime()
	}

	return models.User{
		Id:          parsedUUID,
		Email:       proto_usr.GetEmail(),
		Password:    proto_usr.GetPassword(),
		Role:        proto_usr.GetRole(),
		Nick:        proto_usr.GetNick(),
		Description: proto_usr.GetDescription(),
		Birthday:    birthday,
	}, nil
}

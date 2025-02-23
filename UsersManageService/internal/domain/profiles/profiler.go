package profiles

import (
	"time"
	"usersManageService/internal/domain/models"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UsrToProtoUsr(user models.User) (*umv1.User, error) {
	var birthday *timestamppb.Timestamp
	if !user.Birthday.IsZero() {
		birthday = timestamppb.New(user.Birthday)
	}

	return &umv1.User{
		Id:       user.Id.String(),
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Nick:     user.Nick,
		Birthday: birthday,
	}, nil
}

func ProtoUsrToUsr(proto_usr *umv1.User) (models.User, error) {
	parsedUUID, err := uuid.Parse(proto_usr.GetId())
	if err != nil {
		return models.User{}, err
	}

	var birthday time.Time
	if proto_usr.GetBirthday() != nil {
		birthday = proto_usr.GetBirthday().AsTime()
	}

	return models.User{
		Id:       parsedUUID,
		Email:    proto_usr.GetEmail(),
		Password: proto_usr.GetPassword(),
		Role:     proto_usr.GetRole(),
		Nick:     proto_usr.GetNick(),
		Birthday: birthday,
	}, nil
}

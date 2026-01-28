package grpc_user_info

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	"google.golang.org/grpc"
)

var (
	ErrNilUser = errors.New("converter: user is nil")

	ErrGetUser = errors.New("grpc user info client: failed to get user")
)

type UserInfoService struct {
	client proto.UserServiceClient
}

func NewUserInfoService(conn *grpc.ClientConn) *UserInfoService {
	return &UserInfoService{
		client: proto.NewUserServiceClient(conn),
	}
}

func (u *UserInfoService) GetUser(ctx context.Context, userID uuid.UUID) (_ *entity.User, err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrGetUser, err)
	}()

	resp, err := u.client.GetUser(ctx, &proto.GetUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return protoToEntityUser(resp)
}

func protoToEntityUser(p *proto.User) (*entity.User, error) {
	if p == nil {
		return nil, ErrNilUser
	}
	id, err := uuid.Parse(p.Id)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:               id,
		Name:             p.Name,
		AccessLvl:        entity.AccessLevel(p.AccessLevel),
		RegistrationDate: p.RegistrationDate.AsTime(),
	}, nil
}

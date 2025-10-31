package grpc_user_creator

import (
	"context"

	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/grpc/converter"
	"google.golang.org/grpc"
)

type UserCreatorService struct {
	client proto.UserServiceClient
}

func NewUserCreatorService(conn *grpc.ClientConn) *UserCreatorService {
	return &UserCreatorService{
		client: proto.NewUserServiceClient(conn),
	}
}

func (u *UserCreatorService) CreateUser(ctx context.Context, newUser *entity.User) (*entity.User, error) {
	resp, err := u.client.CreateUser(ctx, &proto.CreateUserRequest{
		Username: newUser.Name,
	})
	if err != nil {
		return nil, err
	}

	user, err := converter.ProtoToEntityUser(resp)
	if err != nil {
		return nil, err
	}

	return user, nil
}

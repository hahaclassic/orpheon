package grpc_user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/grpc/converter"
	"google.golang.org/grpc"
)

type UserService struct {
	client proto.UserServiceClient
}

func NewUserService(conn *grpc.ClientConn) *UserService {
	return &UserService{
		client: proto.NewUserServiceClient(conn),
	}
}

func (u *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	resp, err := u.client.GetUser(ctx, &proto.GetUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return converter.ProtoToDTOUser(resp)
}

func (u *UserService) UpdateUser(ctx context.Context, claims *dto.Claims, req *dto.EditUserRequest) error {
	_, err := u.client.UpdateUser(ctx, &proto.UpdateUserRequest{
		Claims: converter.DTOToProtoClaims(claims),
		User: converter.DTOToProtoUser(&dto.User{
			ID:   req.ID,
			Name: req.Name,
		}),
	})

	return err
}

func (u *UserService) DeleteUser(ctx context.Context, claims *dto.Claims, userID uuid.UUID) error {
	_, err := u.client.DeleteUser(ctx, &proto.DeleteUserRequest{
		Claims: converter.DTOToProtoClaims(claims),
		UserId: userID.String(),
	})

	return err
}

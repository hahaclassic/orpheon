package grpc_ctrl

import (
	"context"

	"github.com/google/uuid"
	proto "github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/usecase"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/providers/grpc/converter"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserController struct {
	proto.UnimplementedUserServiceServer
	service usecase.UserService
}

func NewUserController(userService usecase.UserService) *UserController {
	return &UserController{service: userService}
}

func (s *UserController) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	requestUser := &entity.User{
		Name: req.Username,
	}

	user, err := s.service.CreateUser(ctx, requestUser)
	if err != nil {
		return nil, err
	}

	return converter.EntityToProtoUser(user), nil
}

func (s *UserController) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.User, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	u, err := s.service.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return converter.EntityToProtoUser(u), nil
}

func (s *UserController) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	claims, err := converter.ProtoToEntityClaims(req.Claims)
	if err != nil {
		return nil, err
	}

	u, err := converter.ProtoToEntityUser(req.User)
	if err != nil {
		return nil, err
	}

	if err := s.service.UpdateUser(ctx, claims, u); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserController) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	claims, err := converter.ProtoToEntityClaims(req.Claims)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := s.service.DeleteUser(ctx, claims, userID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

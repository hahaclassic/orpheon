package grpc

import (
	"context"

	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/usecase/user"
	userpb "github.com/hahaclassic/orpheon/services/user-msv/pkg/proto/userpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	userService user.UserService
}

func NewUserGRPCServer(userService user.UserService) *UserGRPCServer {
	return &UserGRPCServer{userService: userService}
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *userpb.UserIDRequest) (*userpb.UserResponse, error) {
	userID, err := entity.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	u, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		User: &userpb.User{
			Id:               u.ID.String(),
			Name:             u.Name,
			Email:            u.Email,
			AccessLevel:      userpb.AccessLevel(u.AccessLvl),
			RegistrationDate: entity.TimestampProto(u.RegistrationDate),
		},
	}, nil
}

func (s *UserGRPCServer) GetMe(ctx context.Context, _ *emptypb.Empty) (*userpb.UserResponse, error) {
	claims := entity.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, user.ErrForbidden
	}

	u, err := s.userService.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		User: &userpb.User{
			Id:               u.ID.String(),
			Name:             u.Name,
			Email:            u.Email,
			AccessLevel:      userpb.AccessLevel(u.AccessLvl),
			RegistrationDate: entity.TimestampProto(u.RegistrationDate),
		},
	}, nil
}

func (s *UserGRPCServer) UpdateMe(ctx context.Context, req *userpb.UpdateUserRequest) (*emptypb.Empty, error) {
	claims := entity.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, user.ErrForbidden
	}

	u := &entity.User{
		ID:    entity.MustParseUUID(req.User.Id),
		Name:  req.User.Name,
		Email: req.User.Email,
		// AccessLvl не обновляем через UpdateMe
	}

	if err := s.userService.UpdateUser(ctx, claims, u); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserGRPCServer) DeleteUser(ctx context.Context, req *userpb.UserIDRequest) (*emptypb.Empty, error) {
	claims := entity.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, user.ErrForbidden
	}

	userID := entity.MustParseUUID(req.Id)
	if err := s.userService.DeleteUser(ctx, claims, userID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

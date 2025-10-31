package grpc_ctrl

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/api/auth-msv/v1/proto"
	jwtproto "github.com/hahaclassic/orpheon/api/jwt/v1/proto"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/usecase"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/grpc/converter"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrEmptyGRPCRequestBody = errors.New("controller: empty grpc request body")
)

type AuthController struct {
	proto.UnimplementedAuthServiceServer
	service usecase.AuthService
}

func NewAuthController(svc usecase.AuthService) *AuthController {
	return &AuthController{service: svc}
}

func (c *AuthController) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.AuthTokens, error) {
	if req == nil {
		return nil, ErrEmptyGRPCRequestBody
	}
	newUser := &entity.User{Name: req.Username}
	creds := &entity.UserCredentials{
		Login:    req.Login,
		Password: req.Password,
	}

	tokens, err := c.service.RegisterUser(ctx, newUser, creds)
	if err != nil {
		return nil, err
	}

	return &proto.AuthTokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (c *AuthController) Login(ctx context.Context, req *proto.UserCredentials) (*proto.AuthTokens, error) {
	if req == nil {
		return nil, ErrEmptyGRPCRequestBody
	}
	creds := &entity.UserCredentials{
		Login:    req.Login,
		Password: req.Password,
	}

	tokens, err := c.service.Login(ctx, creds)
	if err != nil {
		return nil, err
	}

	return &proto.AuthTokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (c *AuthController) Logout(ctx context.Context, refresh *proto.RefreshToken) (*emptypb.Empty, error) {
	if refresh == nil {
		return nil, ErrEmptyGRPCRequestBody
	}
	if err := c.service.Logout(ctx, refresh.Value); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c *AuthController) RefreshTokens(ctx context.Context, refresh *proto.RefreshToken) (*proto.AuthTokens, error) {
	if refresh == nil {
		return nil, ErrEmptyGRPCRequestBody
	}
	tokens, err := c.service.RefreshTokens(ctx, refresh.Value)
	if err != nil {
		return nil, err
	}
	return &proto.AuthTokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (c *AuthController) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, ErrEmptyGRPCRequestBody
	}
	pass := &entity.UserPasswords{
		Old: req.Passwords.Old,
		New: req.Passwords.New,
	}
	claims, err := converter.ProtoToEntityClaims(req.Claims)
	if err != nil {
		return nil, err
	}

	if err := c.service.UpdatePassword(ctx, claims.UserID, pass); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *AuthController) GetClaims(ctx context.Context, access *proto.AccessToken) (*jwtproto.Claims, error) {
	claims, err := c.service.GetClaims(ctx, access.Value)
	if err != nil {
		return nil, err
	}
	return converter.EntityToProtoClaims(claims), nil
}

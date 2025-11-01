package grpc_auth

import (
	"context"

	"github.com/hahaclassic/orpheon/api/auth-msv/v1/proto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/dto"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/grpc/converter"
	"google.golang.org/grpc"
)

type AuthService struct {
	client proto.AuthServiceClient
}

func NewAuthService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{
		client: proto.NewAuthServiceClient(conn),
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, request *dto.RegisterRequest) (*dto.AuthTokens, error) {
	resp, err := a.client.RegisterUser(ctx, &proto.RegisterUserRequest{
		Login:    request.Login,
		Password: request.Password,
		Username: request.Username,
	})
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		Access:  resp.Access,
		Refresh: resp.Refresh,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, creds *dto.UserCredentials) (*dto.AuthTokens, error) {
	resp, err := a.client.Login(ctx, &proto.UserCredentials{
		Login:    creds.Login,
		Password: creds.Password,
	})
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		Access:  resp.Access,
		Refresh: resp.Refresh,
	}, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*dto.AuthTokens, error) {
	resp, err := a.client.RefreshTokens(ctx, &proto.RefreshToken{
		Value: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		Access:  resp.Access,
		Refresh: resp.Refresh,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	_, err := a.client.Logout(ctx, &proto.RefreshToken{
		Value: refreshToken,
	})

	return err
}

func (a *AuthService) GetClaims(ctx context.Context, accessToken string) (*dto.Claims, error) {
	claims, err := a.client.GetClaims(ctx, &proto.AccessToken{
		Value: accessToken,
	})
	if err != nil {
		return nil, err
	}

	return converter.ProtoToDTOClaims(claims)
}

func (a *AuthService) UpdatePassword(ctx context.Context, claims *dto.Claims, passwords *dto.UserPasswords) error {
	_, err := a.client.UpdatePassword(ctx, &proto.UpdatePasswordRequest{
		Claims: converter.DTOToProtoClaims(claims),
		Passwords: &proto.UserPasswords{
			Old: passwords.Old,
			New: passwords.New,
		},
	})

	return err
}

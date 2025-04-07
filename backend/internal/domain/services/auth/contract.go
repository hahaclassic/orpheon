package auth

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, user)  
	Login()
	Logout()
	RefreshToken()
	ValidateTokens()
	UpdatePassword()
}

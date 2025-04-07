package entities

import "github.com/google/uuid"

type Claims struct {
	UserID    uuid.UUID
	AccessLvl AccessLevel
}

type JWTClaims struct {
	UserID    uuid.UUID
	AccessLvl AccessLevel
	Exp       int64
}

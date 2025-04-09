package entity

import "github.com/google/uuid"

type Claims struct {
	UserID    uuid.UUID
	AccessLvl AccessLevel
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID               uuid.UUID
	Name             string
	RegistrationDate time.Time
	BirthDate        time.Time
	AccessLvl        AccessLevel
}

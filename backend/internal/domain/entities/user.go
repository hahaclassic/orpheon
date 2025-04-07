package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID               uuid.UUID
	Name             string
	Description      string
	RegistrationDate time.Time
	BirthDate        time.Time
}

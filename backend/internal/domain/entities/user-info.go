package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID               uuid.UUID
	Name             string
	Status           string
	RegistrationDate time.Time
	BirthDate        time.Time
}

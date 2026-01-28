package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID   `json:"id"`
	Name             string      `json:"name"`
	RegistrationDate time.Time   `json:"registration_date"`
	AccessLvl        AccessLevel `json:"access_lvl"`
}

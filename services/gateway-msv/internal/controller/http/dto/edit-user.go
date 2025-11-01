package dto

import "github.com/google/uuid"

type EditUserRequest struct {
	ID   uuid.UUID
	Name string
}

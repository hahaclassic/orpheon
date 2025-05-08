package entity

import "github.com/google/uuid"

type License struct {
	ID          uuid.UUID
	Title       string
	Description string
}

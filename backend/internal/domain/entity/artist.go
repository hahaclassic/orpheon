package entity

import "github.com/google/uuid"

type ArtistMeta struct {
	ID          uuid.UUID
	Name        string
	Description string
	Country     string
}

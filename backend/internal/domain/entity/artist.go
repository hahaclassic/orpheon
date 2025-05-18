package entity

import "github.com/google/uuid"

type ArtistMeta struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Country     string    `json:"country"`
}

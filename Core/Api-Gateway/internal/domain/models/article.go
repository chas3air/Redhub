package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Tag       string    `json:"tag,omitempty"`
	OwnerId   uuid.UUID `json:"owner_id,omitempty"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	OwnerId   uuid.UUID `json:"owner_id"`
}

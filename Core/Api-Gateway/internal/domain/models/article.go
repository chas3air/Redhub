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
	Tag       string    `json:"tag"`
	OwnerId   uuid.UUID `json:"owner_id"`
}

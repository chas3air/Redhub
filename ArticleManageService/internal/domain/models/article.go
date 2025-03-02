package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

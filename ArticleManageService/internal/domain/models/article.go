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
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	Id        uuid.UUID `json:"id"`
	ArticleId uuid.UUID `json:"article_id"`
	OwnerId   uuid.UUID `json:"owner_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

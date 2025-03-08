package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id,omitempty"`
	ArticleId uuid.UUID `json:"article_id,omitempty"`
	OwnerId   uuid.UUID `json:"owner_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Content   string    `json:"content,omitempty"`
}

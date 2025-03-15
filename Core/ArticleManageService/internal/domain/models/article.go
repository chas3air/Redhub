package models

import (
	"fmt"
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

// String метод для Article
func (a Article) String() string {
	return fmt.Sprintf(
		"Article(ID: %s, CreatedAt: %s, Title: %s, OwnerId: %s)",
		a.Id.String(),
		a.CreatedAt.Format(time.RFC3339),
		a.Title,
		a.OwnerId.String(),
	)
}

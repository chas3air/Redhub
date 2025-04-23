package models

import (
	"fmt"
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

func (a Article) String() string {
	return fmt.Sprintf(
		"Article(ID: %s, CreatedAt: %s, Title: %s, OwnerId: %s)",
		a.Id.String(),
		a.CreatedAt.Format(time.RFC3339),
		a.Title,
		a.OwnerId.String(),
	)
}

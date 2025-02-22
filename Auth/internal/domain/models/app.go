package models

import "github.com/google/uuid"

type App struct {
	Id     uuid.UUID `json:"id,omitempty"`
	Alias  string    `json:"alias"`
	Secret string    `json:"secret"`
}

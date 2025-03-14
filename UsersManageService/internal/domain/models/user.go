package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id",omitempty`
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	Nick        string    `json:"nick"`
	Description string    `json:"description"`
	Birthday    time.Time `json:"birthday"`
}

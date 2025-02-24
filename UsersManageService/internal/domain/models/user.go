package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `gorm:"primaryKey"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Nick     string    `json:"nick"`
	Birthday time.Time `json:"birthday"`
}

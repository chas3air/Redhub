package models

import "time"

type Claims struct {
	Uid  string    `json:"uid"`
	Role string    `json:"role"`
	Exp  time.Time `json:"exp"`
}

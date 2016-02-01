package models

import (
	"time"
)

type UserProfile struct {
	Username string    `json:"username" binding:"required"`
	Fullname string    `json:"fullname" binding:"required"`
	Birthday time.Time `json:"birthday" binding:"required"`
}

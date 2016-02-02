package api

import (
	"time"
)

// UserProfile contains information concerning the user
type User struct {
	// Username is the unique name of a user.
	// There can't be 2 same username on the service.
	Username string `json:"username" binding:"required"`
	// Fullname is the more human readable name containing all the classic
	// first name, family/last name, middle name
	Fullname string `json:"fullname" binding:"required"`
	// Birthday is the birthday date of a user
	Birthday time.Time `json:"birthday" binding:"required"`
}

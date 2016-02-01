package errors

import "errors"

var (
	ErrTokenNotFound     = errors.New("Token not found in sessions cookiestore")
	ErrTokenItemNotFound = errors.New("Item not found in token. Possible attack?")
)

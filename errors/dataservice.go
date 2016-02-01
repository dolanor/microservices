package errors

import (
	"errors"
)

var (
	ErrUnauthorized = errors.New("Not authorized to access this content.")
	ErrDataNotFound = errors.New("This ressource doesn't exist.")

	ErrConnectingEndpoint = errors.New("There was a problem connecting to the endpoint")
)

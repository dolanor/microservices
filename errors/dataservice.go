package errors

import (
	"errors"
)

var (
	// ErrUnauthorized indicates when a user is not authorized to access a content
	ErrUnauthorized = errors.New("Not authorized to access this content.")
	// ErrDataNotFound indicates when a ressource is not available
	ErrDataNotFound = errors.New("This ressource doesn't exist.")
	// ErrConnectingEndpoint indicates when an endpoint doesn't answer, is down,
	// replies unexpectandly
	ErrConnectingEndpoint = errors.New("There was a problem connecting to the endpoint")
)

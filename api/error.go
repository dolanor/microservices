package api

import "errors"

var (
	// ErrTokenNotFound indicates when the Token is not found in the cookie store.
	ErrTokenNotFound = errors.New("Token not found in sessions cookiestore")
	// ErrTokenItemNotFound indicates when an item of the token is not available
	// It could be that somebody tries to forge fake Token if they somehow got
	// the server.SymmetricKey.
	ErrTokenItemNotFound = errors.New("Item not found in token. Possible attack?")

	// ErrUnauthorized indicates when a user is not authorized to access a content
	ErrUnauthorized = errors.New("Not authorized to access this content.")
	// ErrDataNotFound indicates when a ressource is not available
	ErrDataNotFound = errors.New("This ressource doesn't exist.")
	// ErrConnectingEndpoint indicates when an endpoint doesn't answer, is down,
	// replies unexpectandly
	ErrConnectingEndpoint = errors.New("There was a problem connecting to the endpoint")
)

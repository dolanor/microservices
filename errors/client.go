package errors

import "errors"

var (
	// ErrTokenNotFound indicates when the Token is not found in the cookie store.
	ErrTokenNotFound = errors.New("Token not found in sessions cookiestore")
	// ErrTokenItemNotFound indicates when an item of the token is not available
	// It could be that somebody tries to forge fake Token if they somehow got
	// the server.SymmetricKey.
	ErrTokenItemNotFound = errors.New("Item not found in token. Possible attack?")
)

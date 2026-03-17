package auth

import "errors"

var (
	// Errors for the responses.
	ErrTokenNotValid = errors.New("error of access-token: is not valid")
)

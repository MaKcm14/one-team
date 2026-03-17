package auth

import "errors"

var (
	// Errors for the responses.
	ErrBadEncoding       = errors.New("error of encoding the info: wrong encoding was got")
	ErrHandleRequest     = errors.New("error of handling the request")
	ErrInvalidAuthHeader = errors.New("error of authorization header: it wasn't set or invalid")
	ErrInvalidAuthInfo   = errors.New("error of the authorization info: the login or password is wrong")
	ErrTokenNotValid     = errors.New("error of access-token: is not valid")
	ErrWrongAuthInfo     = errors.New("error of the authorization info: it's not full or wrong")
)

package auth

import "errors"

var (
	// Errors for the responses.
	ErrLoginRequired        = errors.New("error of session: it has expired: login required")
	ErrInvalidAuthHeader    = errors.New("error of authorization header: it wasn't set or invalid")
	ErrSignUpUserExists     = errors.New("error of sign up the user: already exists")
	ErrInvalidAuthInfo      = errors.New("error of the authorization info: the login or password is wrong")
	ErrAccessTokenNotValid  = errors.New("error of access-token: is not valid")
	ErrRefreshTokenNotValid = errors.New("error of refresh-token: is not valid")
	ErrWrongAuthInfo        = errors.New("error of the authorization info: it's not full or wrong")

	// Session's storage.
	ErrSessionNotFound    = errors.New("auth: error of searching the session: not found")
	ErrSessionWrongFormat = errors.New("auth: error of session: wrong format")
)

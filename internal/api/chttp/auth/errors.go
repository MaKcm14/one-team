package auth

import "errors"

var (
	// Errors for the responses.
	ErrBadEncoding          = errors.New("error of encoding the info: wrong encoding was got")
	ErrLoginRequired        = errors.New("error of session: it has expired: login required")
	ErrHandleRequest        = errors.New("error of handling the request")
	ErrInvalidAuthHeader    = errors.New("error of authorization header: it wasn't set or invalid")
	ErrRequestInfo          = errors.New("error of the request's info")
	ErrSignUpUserExists     = errors.New("error of sign up the user: already exists")
	ErrInvalidAuthInfo      = errors.New("error of the authorization info: the login or password is wrong")
	ErrAccessTokenNotValid  = errors.New("error of access-token: is not valid")
	ErrRefreshTokenNotValid = errors.New("error of refresh-token: is not valid")
	ErrWrongAuthInfo        = errors.New("error of the authorization info: it's not full or wrong")
)

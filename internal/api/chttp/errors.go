package chttp

import "errors"

var (
	ErrStartController = errors.New("chttp: error of starting the controller")

	ErrBindingScheme = errors.New("error of binding the JSON-scheme")
	ErrRequestData   = errors.New("error value of the request's data")

	ErrUserExists        = errors.New("error of the user's value: it already exists")
	ErrAuthData          = errors.New("error of the authorization data's structure")
	ErrInvalidToken      = errors.New("error of the token: invalid token was got")
	ErrInvalidLoginOrPwd = errors.New("error of the login or the password")
	ErrAuthFailed        = errors.New("error of the authorization process")
)

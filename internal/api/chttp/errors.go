package chttp

import "errors"

var (
	// errors of package's funcs.
	ErrStartController = errors.New("chttp: error of starting the controller")
	ErrInvalidData     = errors.New("chttp: error of data: it's invalid")
	ErrInvalidValue    = errors.New("chttp: error of the value: it's in the invalid format")

	// errors of API responses'.
	ErrAuthData          = errors.New("error of the authorization data's structure")
	ErrAuthFailed        = errors.New("error of the authorization process")
	ErrInvalidLoginOrPwd = errors.New("error of login or password value")
	ErrInvalidToken      = errors.New("error of the token: invalid token was got")
	ErrRequestBody       = errors.New("error of validating the request's body")
	ErrRequestQueryParam = errors.New("error of validating the request's query param")
	ErrServerError       = errors.New("error of processing on the server")
	ErrResourceExists    = errors.New("error of resource: it already exists")
	ErrResourceNotFound  = errors.New("error of resource: it wasn't found")
)

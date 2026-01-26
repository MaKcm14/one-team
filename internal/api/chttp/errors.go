package chttp

import "errors"

var (
	ErrStartController = errors.New("chttp: error of starting the controller")

	ErrBindingScheme = errors.New("error of binding the JSON-scheme")
	ErrRequestData   = errors.New("error value of the request's data")
)

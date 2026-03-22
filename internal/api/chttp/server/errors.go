package server

import "errors"

var (
	// Error for the responses.
	ErrAlreadyExists = errors.New("error of creating the object: already exists")
	ErrBadEncoding   = errors.New("error of encoding the info: wrong encoding was got")
	ErrHandleRequest = errors.New("error of handling the request")
	ErrRequestInfo   = errors.New("error of the request's info")
)

package server

import "errors"

var (
	// Error for the responses.
	ErrBadEncoding   = errors.New("error of encoding the info: wrong encoding was got")
	ErrHandleRequest = errors.New("error of handling the request")
	ErrRequestInfo   = errors.New("error of the request's info")
)

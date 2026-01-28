package jwt

import "errors"

var (
	ErrInvalidToken = errors.New("error of the token containing: it's invalid")
)

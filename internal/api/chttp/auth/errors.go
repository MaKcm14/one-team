package auth

import "errors"

var (
	ErrCredsExtracting = errors.New("auth: error of the credential extracting")
	ErrHashPassword    = errors.New("auth: error of hash the password")
	ErrTokenExtracting = errors.New("auth: error of the token extracting")
)

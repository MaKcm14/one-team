package token

import "errors"

var (
	ErrTokenNotValid  = errors.New("token: error of the token: is not valid")
	ErrTokenVerifying = errors.New("token: error of the token's verifying")
)

package token

import "errors"

var (
	ErrTokenIssue     = errors.New("token: error of issue the token")
	ErrTokenNotValid  = errors.New("token: error of the token: is not valid")
	ErrTokenVerifying = errors.New("token: error of the token's verifying")
)

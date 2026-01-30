package tokens

import "errors"

var (
	ErrInvalidToken = errors.New("tokens: error of the token containing: it's invalid")
	ErrTokenConfig  = errors.New("tokens: error of configuration the token")
	ErrTokenParsing = errors.New("tokens: error of parsing the JWT-token")
)

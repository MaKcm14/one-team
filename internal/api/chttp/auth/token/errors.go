package token

import "errors"

var (
	ErrTokenIssue     = errors.New("token: error of issue the token")
	ErrTokenNotValid  = errors.New("token: error of the token: is not valid")
	ErrTokenVerifying = errors.New("token: error of the token's verifying")

	// Storage's errors.
	ErrRefreshTokenHashNotFound    = errors.New("token: error of the refresh-token: not found")
	ErrRefreshTokenHashWrongFormat = errors.New("token: error of the refresh-token format")
	ErrAccessTokenJTINotFound      = errors.New("token: error of the refresh-token: not found")
	ErrAccessTokenJTIWrongFormat   = errors.New("token: error of the refresh-token format")
)

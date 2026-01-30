package auth

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"auth-train/test/internal/api/chttp/auth/tokens"
)

const (
	basicAuthPrefix  = "Basic "
	bearerAuthPrefix = "Bearer "

	headerAuthorization   = "Authorization"
	headerWWWAuthenticate = "WWW-Authenticate"
)

type Authenticator struct {
	UserToken tokens.UserJWT
}

func NewAuthenticator(conf tokens.AuthJWTConfig) Authenticator {
	return Authenticator{
		UserToken: tokens.NewUserJWT(conf),
	}
}

func ExtractRawToken(ctx echo.Context) (token string, err error) {
	authHeader := ctx.Request().Header.
		Get(headerAuthorization)
	if !strings.HasPrefix(authHeader, bearerAuthPrefix) {
		ctx.Response().Header().
			Set(
				headerWWWAuthenticate,
				fmt.Sprintf("%srealm=\"Restricted\"", bearerAuthPrefix),
			)
		return "", ErrTokenExtracting
	}

	return strings.TrimPrefix(authHeader, bearerAuthPrefix), nil
}

func ExtractCreds(ctx echo.Context) (login, pwd string, err error) {
	authHeader := ctx.Request().Header.
		Get(headerAuthorization)
	if !strings.HasPrefix(authHeader, basicAuthPrefix) {
		ctx.Response().Header().
			Set(
				headerWWWAuthenticate,
				fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix),
			)
		return "", "", ErrCredsExtracting
	}

	payload, err := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(authHeader, basicAuthPrefix),
	)
	if err != nil {
		ctx.Response().Header().
			Set(
				headerWWWAuthenticate,
				fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix),
			)
		return "", "", ErrCredsExtracting
	}

	creds := strings.SplitN(string(payload), ":", 2)
	if len(creds) != 2 {
		ctx.Response().Header().
			Set(
				headerWWWAuthenticate,
				fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix),
			)
		return "", "", ErrCredsExtracting
	}
	return creds[0], creds[1], nil
}

func HashPassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

func IsPasswordEqual(hashPwd []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashPwd, password)
}

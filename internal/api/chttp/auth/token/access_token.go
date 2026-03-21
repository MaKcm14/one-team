package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

const (
	IssuerName = "hrm.oneteam.com"

	AccessTokenTTL = 5 * time.Minute
)

type Claims struct {
	jwt.RegisteredClaims

	SessionID string      `json:"sid"`
	UserData  user.Claims `json:"user"`
}

type AccessToken struct {
	cfg config.AuthConfig
}

func NewAccessToken(cfg config.AuthConfig) AccessToken {
	return AccessToken{
		cfg: cfg,
	}
}

func (a AccessToken) getKey(_ *jwt.Token) (any, error) {
	return []byte(a.cfg.Secret), nil
}

func (a AccessToken) VerifyAccessToken(accessToken string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		a.getKey,
		jwt.WithIssuer(IssuerName),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %s", ErrTokenVerifying, err)
	}

	if !token.Valid {
		return Claims{}, ErrTokenNotValid
	}
	return claims, nil
}

func (a AccessToken) IssueAccessToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(a.cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrTokenIssue, err)
	}
	return signedToken, nil
}

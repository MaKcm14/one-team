package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth-train/test/internal/entity"
)

const (
	issuer = "auth.train.com"
)

type UserJWT struct {
	secret []byte
	method jwt.SigningMethod
}

func NewUserJWT(jwtConfig AuthJWTConfig) UserJWT {
	return UserJWT{
		secret: jwtConfig.Secret,
		method: jwtConfig.Method,
	}
}

func (u UserJWT) NewUserToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"iss": issuer,
		"sub": user.ID,
		"aud": issuer,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(
		u.method,
		claims,
	).SignedString(u.secret)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrTokenConfig, err)
	}
	return token, nil
}

func (u UserJWT) keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (any, error) {
		return u.secret, nil
	}
}

func (u UserJWT) VerifyUserJWT(accessToken string) (entity.UserID, error) {
	token, err := jwt.Parse(
		accessToken,
		u.keyFunc(),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{u.method.Alg()}),
	)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	} else if !token.Valid {
		return 0, ErrTokenParsing
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return entity.UserID(claims["sub"].(float64)), nil
}

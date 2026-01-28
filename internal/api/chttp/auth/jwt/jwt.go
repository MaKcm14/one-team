package jwt

import (
	"auth-train/test/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer = "auth.train.com"
)

type UserJWT struct {
	HMACSecret string
}

func NewUserJWT(secret string) UserJWT {
	return UserJWT{
		HMACSecret: secret,
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

	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString([]byte(u.HMACSecret))
}

func (u UserJWT) keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (any, error) {
		return []byte(u.HMACSecret), nil
	}
}

func (u UserJWT) VerifyUserJWT(accessToken string) (entity.UserID, error) {
	token, err := jwt.Parse(
		accessToken,
		u.keyFunc(),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return entity.UserID(claims["sub"].(float64)), nil
}

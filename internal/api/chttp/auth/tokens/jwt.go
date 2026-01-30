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
	claims := claimsJWT{
		UserClaims: userToUserClaims(user),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  []string{issuer},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
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

func (u UserJWT) VerifyUserJWT(accessToken string) (UserClaims, error) {
	claims := claimsJWT{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		u.keyFunc(),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{u.method.Alg()}),
	)

	if err != nil {
		return UserClaims{}, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	} else if !token.Valid {
		return UserClaims{}, ErrTokenParsing
	}
	return claims.UserClaims, nil
}

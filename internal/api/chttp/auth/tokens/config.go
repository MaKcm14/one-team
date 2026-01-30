package tokens

import "github.com/golang-jwt/jwt/v5"

type AuthJWTConfig struct {
	Secret []byte
	Method jwt.SigningMethod
}

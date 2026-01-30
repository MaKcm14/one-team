package tokens

import (
	"auth-train/test/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type claimsJWT struct {
	UserClaims UserClaims `json:"sub"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	ID          entity.UserID
	AdminStatus bool
}

func userToUserClaims(user entity.User) UserClaims {
	return UserClaims{
		ID:          user.ID,
		AdminStatus: user.Profile.AdminStatus,
	}
}

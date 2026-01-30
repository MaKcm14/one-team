package tokens

import (
	"github.com/golang-jwt/jwt/v5"

	"auth-train/test/internal/entity"
)

type UserClaimsJWT struct {
	UserPayload UserPayloadJWT `json:"sub"`
	jwt.RegisteredClaims
}

type UserPayloadJWT struct {
	ID          entity.UserID `json:"user_id"`
	AdminStatus bool          `json:"admin"`
}

func userToUserPayloadJWT(user entity.User) UserPayloadJWT {
	return UserPayloadJWT{
		ID:          user.ID,
		AdminStatus: user.Profile.AdminStatus,
	}
}

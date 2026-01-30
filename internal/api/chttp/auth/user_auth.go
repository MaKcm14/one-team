package auth

import (
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/entity"
)

type UserAuthenticator struct {
	UserToken tokens.UserJWT
	validJTI  map[entity.UserID]string
}

func NewUserAuthenticator(conf tokens.AuthJWTConfig) UserAuthenticator {
	return UserAuthenticator{
		UserToken: tokens.NewUserJWT(conf),
		validJTI:  make(map[entity.UserID]string, 100),
	}
}

func (u *UserAuthenticator) IsTokenCoolDown(userID entity.UserID, tokenID string) bool {
	if id, ok := u.validJTI[userID]; ok && id == tokenID {
		return false
	}
	return true
}

func (u *UserAuthenticator) TokenCoolDown(userID entity.UserID) {
	delete(u.validJTI, userID)
}

func (u *UserAuthenticator) RegisterToken(userID entity.UserID, tokenID string) {
	u.validJTI[userID] = tokenID
}

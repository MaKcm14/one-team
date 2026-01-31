package auth

import (
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/entity"
	"sync"
)

type UserAuthenticator struct {
	UserToken tokens.UserJWT
	validJTI  map[entity.UserID]string
	mx        sync.RWMutex
}

func NewUserAuthenticator(conf tokens.AuthJWTConfig) UserAuthenticator {
	return UserAuthenticator{
		UserToken: tokens.NewUserJWT(conf),
		validJTI:  make(map[entity.UserID]string, 100),
	}
}

func (u *UserAuthenticator) IsTokenCoolDown(userID entity.UserID, tokenID string) bool {
	u.mx.RLock()
	defer u.mx.RUnlock()

	if id, ok := u.validJTI[userID]; ok && id == tokenID {
		return false
	}
	return true
}

func (u *UserAuthenticator) TokenCoolDown(userID entity.UserID) {
	u.mx.Lock()
	delete(u.validJTI, userID)
	u.mx.Unlock()
}

func (u *UserAuthenticator) RegisterToken(userID entity.UserID, tokenID string) {
	u.mx.Lock()
	u.validJTI[userID] = tokenID
	u.mx.Unlock()
}

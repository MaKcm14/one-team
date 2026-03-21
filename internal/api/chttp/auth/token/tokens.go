package token

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// TokenStorage stores the access- and refresh- tokens:
// KEY: session_id; VALUE: JTI
// KEY: session_id; VALUE: Hash(Refresh)
type TokenStorage struct {
	AccessTokens  *cache.Cache
	RefreshTokens *cache.Cache
}

func NewTokenStorage() TokenStorage {
	return TokenStorage{
		AccessTokens:  cache.New(AccessTokenTTL, 30*time.Minute),
		RefreshTokens: cache.New(RefreshTokenTTL, 30*time.Minute),
	}
}

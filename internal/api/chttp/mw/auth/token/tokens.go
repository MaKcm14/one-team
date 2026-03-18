package token

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type TokenStorage struct {
	AccessTokens  *cache.Cache
	RefreshTokens *cache.Cache
}

func NewTokenStorage() TokenStorage {
	return TokenStorage{
		AccessTokens:  cache.New(AccessTokenTTL, 5*time.Minute),
		RefreshTokens: cache.New(RefreshTokenTTL, 30*time.Minute),
	}
}

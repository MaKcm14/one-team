package token

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type TokenStorage struct {
	RefreshTokens *cache.Cache
}

func NewTokenStorage() TokenStorage {
	return TokenStorage{
		RefreshTokens: cache.New(RefreshTokenTTL, 30*time.Minute),
	}
}

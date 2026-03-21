package token

import (
	"fmt"
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

func (t TokenStorage) GetHashRefreshToken(sid string) (string, error) {
	val, ok := t.RefreshTokens.Get(sid)
	if !ok {
		return "", fmt.Errorf("%w: has expired or wasn't existed", ErrRefreshTokenHashNotFound)
	}

	hashRefreshToken, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%w: token's converting format is wrong", ErrRefreshTokenHashWrongFormat)
	}
	return hashRefreshToken, nil
}

func (t TokenStorage) SetHashRefreshToken(sid string, hash string, ttl time.Duration) {
	t.RefreshTokens.Set(sid, hash, ttl)
}

func (t TokenStorage) GetAccessTokenJTI(sid string) (string, error) {
	val, ok := t.AccessTokens.Get(sid)
	if !ok {
		return "", fmt.Errorf("%w: has expired or wasn't existed", ErrAccessTokenJTINotFound)
	}

	jti, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%w: token's converting format is wrong", ErrAccessTokenJTIWrongFormat)
	}
	return jti, nil
}

func (t TokenStorage) SetAccessTokenJTI(sid string, jti string, ttl time.Duration) {
	t.RefreshTokens.Set(sid, jti, ttl)
}

package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	entity "github.com/MaKcm14/one-team/internal/entity/user"
)

const sessionIDCookieKey = "session_id"

type UserSession struct {
	User entity.User
	Role entity.Role
}

type SessionConfig struct {
	Sessions *cache.Cache
	Writer   sessions.CookieStore
}

func NewSessionConfig() SessionConfig {
	return SessionConfig{
		Sessions: cache.New(24*time.Hour, time.Hour),
		Writer:   *sessions.NewCookieStore(securecookie.GenerateRandomKey(32)),
	}
}

func (a Authenticator) createSession() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	a.session.Sessions.Set(id.String(), UserSession{}, token.AccessTokenTTL)
	return id.String(), nil
}

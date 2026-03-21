package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

const sessionIDCookieKey = "session_id"

type SessionConfig struct {
	Sessions *cache.Cache
	Writer   sessions.CookieStore
}

func NewSessionConfig(cfg config.AuthConfig) SessionConfig {
	return SessionConfig{
		Sessions: cache.New(24*time.Hour, time.Hour),
		Writer:   *sessions.NewCookieStore([]byte(cfg.SessionKey)),
	}
}

func (a Authenticator) createSession() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	a.session.Sessions.Set(id.String(), user.UserSession{}, token.AccessTokenTTL)
	return id.String(), nil
}

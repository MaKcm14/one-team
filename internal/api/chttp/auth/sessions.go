package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

const (
	sessionIDCookieKey = "session_id"
	SessionIDCtxKey    = "user-session-id"
)

type SessionConfig struct {
	Sessions *cache.Cache
	Writer   *sessions.CookieStore
}

func NewSessionConfig(cfg config.AuthConfig) SessionConfig {
	return SessionConfig{
		Sessions: cache.New(24*time.Hour, time.Hour),
		Writer:   sessions.NewCookieStore([]byte(cfg.SessionKey)),
	}
}

func (s SessionConfig) GetSession(sid string) (user.UserSession, error) {
	val, ok := s.Sessions.Get(sid)
	if !ok {
		return user.UserSession{}, fmt.Errorf("%w: session has expired or wasn't existed", ErrSessionNotFound)
	}

	session, ok := val.(user.UserSession)
	if !ok {
		return user.UserSession{}, fmt.Errorf("%w: other format expected for converting", ErrSessionWrongFormat)
	}
	return session, nil
}

func (s SessionConfig) Set(sid string, session user.UserSession, ttl time.Duration) {
	s.Sessions.Set(sid, session, ttl)
}

func (a Authenticator) createSession() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	a.session.Sessions.Set(id.String(), user.UserSession{}, token.AccessTokenTTL)
	return id.String(), nil
}

func ExtractSessionIDFromCtx(ctx echo.Context) (string, error) {
	val := ctx.Get(SessionIDCtxKey)
	sessionID, ok := val.(string)
	if !ok {
		return "", errors.New("session_id wasn't set for the response")
	}
	return sessionID, nil
}

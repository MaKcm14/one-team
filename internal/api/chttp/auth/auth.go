package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type errorResponse struct {
	Error string `json:"error"`
}

type Authenticator struct {
	log         *slog.Logger
	acToken     token.AccessToken
	refToken    token.RefreshToken
	tokens      token.TokenStorage
	session     SessionConfig
	authService user.IAuthService
}

func NewMW(
	log *slog.Logger,
	cfg config.AuthConfig,
	authService user.IAuthService,
) Authenticator {
	return Authenticator{
		log:         log,
		acToken:     token.NewAccessToken(cfg),
		refToken:    token.NewRefreshToken(cfg),
		authService: authService,
		tokens:      token.NewTokenStorage(),
		session:     NewSessionConfig(cfg),
	}
}

func (a Authenticator) VerifyAccessTokenMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const bearerAuth = "Bearer "

			authHeader := ctx.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, bearerAuth) {
				a.log.Warn(fmt.Sprintf("Warn of parsing the 'Authorization' header while verify access-token"))
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrAccessTokenNotValid.Error(),
				})
			}

			session, err := a.session.Writer.Get(ctx.Request(), sessionIDCookieKey)
			if err != nil {
				a.log.Error(fmt.Sprintf("Error of getting the session from request: %s", err))
				return ctx.JSON(http.StatusBadRequest, errorResponse{
					Error: ErrRequestInfo.Error(),
				})
			}
			rawSessionID := session.Values[sessionIDCookieKey]

			sessionID, ok := rawSessionID.(string)
			if !ok {
				a.log.Error(fmt.Sprintf("Error of getting and converting the session from the given format: %s", err))
				return ctx.JSON(http.StatusBadRequest, errorResponse{
					Error: ErrRequestInfo.Error(),
				})
			}

			_, expAt, ok := a.tokens.AccessTokens.GetWithExpiration(sessionID)
			if !ok || expAt.Before(time.Now()) {
				a.log.Warn(fmt.Sprintf("Warn of checking the token: it's expired"))
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrAccessTokenNotValid.Error(),
				})
			}

			claims, err := a.acToken.VerifyAccessToken(authHeader[len(bearerAuth):])
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrAccessTokenNotValid.Error(),
				})
			}

			ctx.Set(TokenClaimsCtxKey, claims)

			return next(ctx)
		}
	}
}

func (a Authenticator) DebugPrintCaches() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)

			acTokenItems := a.tokens.AccessTokens.Items()
			refTokenItems := a.tokens.RefreshTokens.Items()
			sessionItems := a.session.Sessions.Items()

			a.log.Debug("ACCESS_TOKEN_ITEMS_CACHE")
			for key, item := range acTokenItems {
				a.log.Debug(fmt.Sprintf("%s %s", key, item.Object))
			}

			a.log.Debug("REFRESH_TOKEN_ITEMS_CACHE")
			for key, item := range refTokenItems {
				a.log.Debug(fmt.Sprintf("%s %s", key, item.Object))
			}

			a.log.Debug("SESSIONS_ITEMS_CACHE")
			for key, item := range sessionItems {
				a.log.Debug(fmt.Sprintf("%s %s", key, item.Object))
			}
			return err
		}
	}
}

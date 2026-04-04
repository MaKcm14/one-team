package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Authenticator struct {
	log         *slog.Logger
	acToken     token.AccessToken
	refToken    token.RefreshToken
	tokens      token.TokenStorage
	session     SessionConfig
	authService user.IAuthService
	isSysInit   bool
}

func NewAuthenticator(
	log *slog.Logger,
	cfg config.AuthConfig,
	session SessionConfig,
	tokenStorage token.TokenStorage,
	authService user.IAuthService,
) Authenticator {
	return Authenticator{
		log:         log,
		acToken:     token.NewAccessToken(cfg),
		refToken:    token.NewRefreshToken(cfg),
		authService: authService,
		tokens:      tokenStorage,
		session:     session,
	}
}

func (a Authenticator) VerifyAccessTokenMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const bearerAuth = "Bearer "

			authHeader := ctx.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, bearerAuth) {
				a.log.Warn(fmt.Sprintf("Warn of parsing the 'Authorization' header while verify access-token"))
				return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
					Error: ErrAccessTokenNotValid.Error(),
				})
			}

			claims, err := a.acToken.VerifyAccessToken(authHeader[len(bearerAuth):])
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
					Error: ErrAccessTokenNotValid.Error(),
				})
			}
			ctx.Set(TokenClaimsCtxKey, claims)

			return next(ctx)
		}
	}
}

func ExtractClaimsFromCtx(ctx echo.Context) (token.Claims, error) {
	val := ctx.Get(TokenClaimsCtxKey)
	claims, ok := val.(token.Claims)
	if !ok {
		return token.Claims{}, errors.New("claims wasn't set in context")
	}
	return claims, nil
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

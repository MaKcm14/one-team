package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/mw/auth/token"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type errorResponse struct {
	Error string `json:"error"`
}

type Authenticator struct {
	acToken     token.AccessToken
	refToken    token.RefreshToken
	tokens      token.TokenStorage
	session     SessionConfig
	authService user.IAuthService
}

func NewMW(
	cfg config.AuthConfig,
	authService user.IAuthService,
) Authenticator {
	return Authenticator{
		acToken:     token.NewAccessToken(cfg),
		authService: authService,
		tokens:      token.NewTokenStorage(),
	}
}

func (a Authenticator) VerifyAccessTokenMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const bearerAuth = "Bearer "

			authHeader := ctx.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, bearerAuth) {
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrTokenNotValid.Error(),
				})
			}

			session, err := a.session.Writer.Get(ctx.Request(), sessionIDCookieKey)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, errorResponse{
					Error: ErrHandleRequest.Error(),
				})
			}
			rawSessionID := session.Values[sessionIDCookieKey]

			sessionID, ok := rawSessionID.(string)
			if !ok {
				return ctx.JSON(http.StatusInternalServerError, errorResponse{
					Error: ErrHandleRequest.Error(),
				})
			}
			_, expAt, ok := a.tokens.AccessTokens.GetWithExpiration(sessionID)
			if !ok || expAt.Before(time.Now()) {
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrTokenNotValid.Error(),
				})
			}

			claims, err := a.acToken.VerifyAccessToken(authHeader[len(bearerAuth):])
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrTokenNotValid.Error(),
				})
			}

			ctx.Set(TokenClaimsCtxKey, claims)

			return next(ctx)
		}
	}
}

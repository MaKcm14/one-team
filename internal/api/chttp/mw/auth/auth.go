package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/mw/auth/token"
	"github.com/MaKcm14/one-team/internal/config"
)

type errorResponse struct {
	Error string `json:"error"`
}

type AuthMiddleware struct {
	acToken token.AccessToken
}

func NewMW(cfg config.AuthConfig) AuthMiddleware {
	return AuthMiddleware{
		acToken: token.NewAccessToken(cfg),
	}
}

func (a AuthMiddleware) VerifyAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const bearerAuth = "Bearer "

			authHeader := ctx.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, bearerAuth) {
				return ctx.JSON(http.StatusUnauthorized, errorResponse{
					Error: ErrTokenNotValid.Error(),
				})
			}
			rawToken := strings.TrimPrefix(authHeader, bearerAuth)

			claims, err := a.acToken.VerifyAccessToken(rawToken)
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

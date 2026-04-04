package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
)

func (a Authenticator) GrantAdminAccessMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			claims, _ := ExtractClaimsFromCtx(ctx)

			if claims.UserData.Role != entity.AdminRole {
				a.log.Error(fmt.Sprintf("Try to access to admin's resource from '%s'",
					logUnauthorizedEvent(claims)))
				return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
					Error: ErrPermissionDenied.Error(),
				})
			}
			return next(ctx)
		}
	}
}

func (a Authenticator) GrantAdminOrHRManagerAccessMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			claims, _ := ExtractClaimsFromCtx(ctx)

			if role := claims.UserData.Role; role != entity.AdminRole && role != entity.HRManagerRole {
				a.log.Error(fmt.Sprintf("Try to access to admin/hrmanager's resource from %s",
					logUnauthorizedEvent(claims)))
				return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
					Error: ErrPermissionDenied.Error(),
				})
			}
			return next(ctx)
		}
	}
}

func (a Authenticator) GrantAllAccessMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			claims, _ := ExtractClaimsFromCtx(ctx)

			if role := claims.UserData.Role; role != entity.AdminRole && role != entity.HRManagerRole && role != entity.AnalystRole {
				a.log.Error(fmt.Sprintf("Try to access to admin/hrmanager/analyst's resource from %s",
					logUnauthorizedEvent(claims)))
				return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
					Error: ErrPermissionDenied.Error(),
				})
			}
			return next(ctx)
		}
	}
}

func logUnauthorizedEvent(claims token.Claims) string {
	return fmt.Sprintf("session_id=%s; login=%s; role=%s",
		claims.SessionID, claims.UserData.Login, claims.UserData.Role)
}

package admin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/services/usecase/root"
	"github.com/labstack/echo/v4"
)

type AdminRouter struct {
	log     *slog.Logger
	tokens  token.TokenStorage
	session auth.SessionConfig

	rootService root.IRootService
}

func NewAdminRouter(
	log *slog.Logger,
	tokens token.TokenStorage,
	session auth.SessionConfig,
	rootService root.IRootService,
) AdminRouter {
	return AdminRouter{
		log:         log,
		tokens:      tokens,
		session:     session,
		rootService: rootService,
	}
}

type userResponse struct {
	User             root.UserDTO `json:"user"`
	JTI              string       `json:"jti"`
	HashRefreshToken string       `json:"hash_refresh_token"`
	SessionID        string       `json:"session_id"`
}

func (a AdminRouter) HandlerAdminGetUsers(eCtx echo.Context) error {
	type response struct {
		Users []userResponse `json:"users"`
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	list, err := a.rootService.GetUsers(ctx)
	if err != nil {
		a.log.Error(fmt.Sprintf("Error of getting the users: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}

	users := make([]userResponse, 0, len(list))
	for _, item := range list {
		sid, _ := a.session.GetSIDForLogin(item.Login)
		token, _ := a.tokens.GetAccessTokenJTI(sid)
		hashRefresh, _ := a.tokens.GetHashRefreshToken(sid)

		users = append(users, userResponse{
			User:             item,
			JTI:              token,
			SessionID:        sid,
			HashRefreshToken: hashRefresh,
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Users: users,
	})
}

func (a AdminRouter) HandlerAdminGetRoles(eCtx echo.Context) error {
	type response struct {
		Roles []root.Role `json:"roles"`
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	list, err := a.rootService.GetRoles(ctx)
	if err != nil {
		a.log.Error(fmt.Sprintf("Error of getting the roles: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.JSON(http.StatusOK, response{
		Roles: list,
	})
}

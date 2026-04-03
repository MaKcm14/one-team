package admin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	entity "github.com/MaKcm14/one-team/internal/entity/user"
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

func (a AdminRouter) HandlerAdminSessionFlush(eCtx echo.Context) error {
	delSessionID, err := server.ValidateSessionID(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	delSession, err := a.session.GetSession(delSessionID)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: session has expired or wasn't exist", server.ErrRequestInfo),
		})
	}

	if delSession.UserClaims.Role == entity.AdminRole {
		a.log.Error(fmt.Sprintf("Try to flush the admin's session"))
		return eCtx.JSON(http.StatusUnauthorized, server.ErrorResponse{
			Error: "unable to flush the admin's session",
		})
	}
	a.session.Sessions.Delete(delSessionID)
	a.session.Sessions.Delete(delSession.UserClaims.Login)
	a.tokens.AccessTokens.Delete(delSessionID)
	a.tokens.RefreshTokens.Delete(delSessionID)

	return eCtx.NoContent(http.StatusOK)
}

func (a AdminRouter) HandlerAdminDeleteUser(eCtx echo.Context) error {
	login, err := server.ValidateLogin(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: fmt.Sprintf("%s: %s", server.ErrRequestInfo, err),
		})
	}

	claims, err := auth.ExtractClaimsFromCtx(eCtx)
	if err != nil {
		a.log.Error("Error of extracting the user's claims from context")
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}

	if claims.UserData.Login == login {
		a.log.Error("Try to delete the user with the login equals to the deleter")
		return eCtx.JSON(http.StatusUnauthorized, server.ErrorResponse{
			Error: "unable to delete the current user",
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 5*time.Second)
	defer cancel()

	err = a.rootService.DeleteUser(ctx, login)
	if err != nil {
		if errors.Is(err, root.ErrUserNotFound) {
			a.log.Warn("Try to delete the unexisting user")
			return eCtx.JSON(http.StatusNotFound, server.ErrorResponse{
				Error: fmt.Sprintf("%s: unable to delete the unexisting user", server.ErrRequestInfo),
			})
		} else if errors.Is(err, root.ErrUnableToDeleteAdmin) {
			a.log.Error("Error of deleting the admin: the operation is restricted")
			return eCtx.JSON(http.StatusUnauthorized, server.ErrorResponse{
				Error: fmt.Sprintf("%s: unable to delete the admin: the operation is restricted", server.ErrRequestInfo),
			})
		}
		a.log.Error(fmt.Sprintf("Error occured while trying to delete the user: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusOK)
}

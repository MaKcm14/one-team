package chttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"

	"auth-train/test/internal/api/chttp/auth"
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

const (
	userJWTCtxKey = "userJWTClaims"
)

func (h *HttpController) handlerSignUp(ctx echo.Context) error {
	const op = "chttp.handlerSignUp"

	userCfg := repo.UserConfig{}
	if err := ctx.Bind(&userCfg); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestBody.Error(),
		})
	}

	err := validateUser(
		userCfg,
		withNameValidator(),
		withSurnameValidator(),
		withPassportValidator(),
	)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestBody.Error(),
		})
	}

	if _, err := h.store.GetUserByPassport(userCfg.Passport); err == nil {
		return ctx.JSON(http.StatusConflict, errorResponse{
			fmt.Sprintf("error of sign-up: %s", ErrResourceExists),
		})
	}

	login, pwd, err := auth.ExtractCreds(ctx)
	if login != userCfg.Passport {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrAuthData.Error(),
		})
	} else if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrAuthData.Error(),
		})
	}

	userCfg.PwdHash, err = auth.HashPassword(pwd)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			fmt.Sprintf("error of sign-up procedure: %s", ErrServerError),
		})
	}
	return ctx.JSON(http.StatusCreated, h.store.CreateUser(userCfg))
}

func (h *HttpController) handlerLogin(ctx echo.Context) error {
	const op = "chttp.handlerLogin"
	type authInfo struct {
		AccessToken string `json:"token"`
	}

	login, pwd, err := auth.ExtractCreds(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrAuthData.Error(),
		})
	}

	user, err := h.store.GetUserByPassport(login)
	if err != nil {
		return ctx.JSON(http.StatusConflict, errorResponse{
			ErrInvalidLoginOrPwd.Error(),
		})
	}

	if err := auth.IsPasswordEqual(user.Profile.PwdHash, []byte(pwd)); err != nil {
		return ctx.JSON(http.StatusConflict, errorResponse{
			ErrInvalidLoginOrPwd.Error(),
		})
	}

	tokenJWT, err := h.auth.UserAuth.UserToken.NewUserToken(user)
	if err != nil {
		h.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			ErrAuthFailed.Error(),
		})
	}
	h.auth.UserAuth.RegisterToken(user.ID, tokenJWT.JTI)
	return ctx.JSON(http.StatusOK, authInfo{tokenJWT.Token})
}

func (h *HttpController) handlerVerifyToken(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := auth.ExtractRawToken(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrAuthData.Error(),
			})
		}

		claims, err := h.auth.UserAuth.UserToken.VerifyUserJWT(token)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}

		_, err = h.store.GetUser(claims.UserPayload.ID)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}

		if h.auth.UserAuth.IsTokenCoolDown(claims.UserPayload.ID, claims.ID) {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}
		ctx.Set(userJWTCtxKey, claims)
		return handler(ctx)
	}
}

func (h *HttpController) handlerGetUserAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, ok := ctx.Get(userJWTCtxKey).(tokens.UserClaimsJWT)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}

		id, err := validateUserID(ctx.QueryParam(userIDParamName))
		if !claims.UserPayload.AdminStatus && err == nil && claims.UserPayload.ID != id {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		} else if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrRequestQueryParam.Error(),
			})
		}
		return handler(ctx)
	}
}

func (h *HttpController) handlerGetUserListAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, ok := ctx.Get(userJWTCtxKey).(tokens.UserClaimsJWT)
		if !ok || !claims.UserPayload.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		return handler(ctx)
	}
}

func (h *HttpController) handlerDeleteUserAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, ok := ctx.Get(userJWTCtxKey).(tokens.UserClaimsJWT)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}

		id, err := validateUserID(ctx.QueryParam(userIDParamName))
		if !claims.UserPayload.AdminStatus && err == nil && claims.UserPayload.ID != id {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		} else if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrRequestQueryParam.Error(),
			})
		}
		return handler(ctx)
	}
}

func (h *HttpController) handlerSetMoneyAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		type request struct {
			ID    entity.UserID `json:"user_id"`
			Money float64       `json:"money"`
		}
		var req request

		claims, ok := ctx.Get(userJWTCtxKey).(tokens.UserClaimsJWT)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}

		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrRequestBody.Error(),
			})
		}

		if req.ID != claims.UserPayload.ID && !claims.UserPayload.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		raw, _ := json.Marshal(req)
		ctx.Request().Body = io.NopCloser(bytes.NewReader(raw))

		return handler(ctx)
	}
}

func (h *HttpController) handlerSetAdminStatusAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		claims, ok := ctx.Get(userJWTCtxKey).(tokens.UserClaimsJWT)
		if !ok || !claims.UserPayload.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		return handler(ctx)
	}
}

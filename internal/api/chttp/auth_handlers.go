package chttp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"auth-train/test/internal/api/chttp/auth"
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

const (
	userClaimsKey = "user_claims"
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

	token, err := h.auth.UserToken.NewUserToken(user)
	if err != nil {
		h.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			ErrAuthFailed.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, authInfo{token})
}

func (h *HttpController) handlerVerifyToken(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := auth.ExtractRawToken(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrAuthData.Error(),
			})
		}

		userClaims, err := h.auth.UserToken.VerifyUserJWT(token)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}

		_, err = h.store.GetUser(userClaims.ID)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}
		ctx.Set(userClaimsKey, userClaims)
		return handler(ctx)
	}
}

func (h *HttpController) handlerGetUserAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userClaims, ok := ctx.Get(userClaimsKey).(tokens.UserClaims)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}

		id, err := validateUserID(ctx.QueryParam(userIDParamName))
		if !userClaims.AdminStatus && err == nil && userClaims.ID != id {
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
		userClaims, ok := ctx.Get(userClaimsKey).(tokens.UserClaims)
		if !ok || !userClaims.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		return handler(ctx)
	}
}

func (h *HttpController) handlerDeleteUserAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userClaims, ok := ctx.Get(userClaimsKey).(tokens.UserClaims)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}

		id, err := validateUserID(ctx.QueryParam(userIDParamName))
		if !userClaims.AdminStatus && err == nil && userClaims.ID != id {
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

		userClaims, ok := ctx.Get(userClaimsKey).(tokens.UserClaims)
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

		if req.ID != userClaims.ID && !userClaims.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		return handler(ctx)
	}
}

func (h *HttpController) handlerSetAdminStatusAC(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userClaims, ok := ctx.Get(userClaimsKey).(tokens.UserClaims)
		if !ok || !userClaims.AdminStatus {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrPermissionDenied.Error(),
			})
		}
		return handler(ctx)
	}
}

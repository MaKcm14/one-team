package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"
	"github.com/labstack/echo/v4"
)

func (a Authenticator) HandlerLogin(eCtx echo.Context) error {
	const basicAuth = "Basic "

	authHeader := eCtx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, basicAuth) {
		return eCtx.JSON(http.StatusUnauthorized, errorResponse{
			ErrInvalidAuthHeader.Error(),
		})
	}

	rawCreds, err := base64.StdEncoding.DecodeString(authHeader[len(basicAuth):])
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, errorResponse{
			Error: ErrBadEncoding.Error(),
		})
	}

	creds := strings.Split(string(rawCreds), ":")
	if len(creds) != 2 {
		return eCtx.JSON(http.StatusBadRequest, errorResponse{
			Error: ErrWrongAuthInfo.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	err = a.authService.Login(ctx, user.Credentials{
		Login:    creds[0],
		Password: creds[1],
	})
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrWrongPassword) {
			return eCtx.JSON(http.StatusUnauthorized, errorResponse{
				Error: ErrInvalidAuthInfo.Error(),
			})
		}
		return eCtx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}

	// TODO: add here creating the access- and refresh- token for the user.
	return nil
}

package chttp

import (
	"auth-train/test/internal/api/chttp/auth/jwt"
	"auth-train/test/internal/repo"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

const (
	basicAuthPrefix  = "Basic "
	bearerAuthPrefix = "Bearer "
)

type authenticator struct {
	userToken jwt.UserJWT
}

func newAuthenticator(secret string) authenticator {
	return authenticator{
		userToken: jwt.NewUserJWT(secret),
	}
}

func (a authenticator) extractRawToken(ctx echo.Context) (token string, err error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, bearerAuthPrefix) {
		ctx.Response().Header().Set("WWW-Authenticate",
			fmt.Sprintf("%srealm=\"Restricted\"", bearerAuthPrefix))
		return "", ErrAuthData
	}

	return strings.TrimPrefix(authHeader, bearerAuthPrefix), nil
}

func (a authenticator) extractCreds(ctx echo.Context) (login, pwd string, err error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, basicAuthPrefix) {
		ctx.Response().Header().Set("WWW-Authenticate",
			fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix))
		return "", "", ErrAuthData
	}

	payload, err := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(authHeader, basicAuthPrefix),
	)
	if err != nil {
		ctx.Response().Header().Set("WWW-Authenticate",
			fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix))
		return "", "", ErrAuthData
	}

	creds := strings.SplitN(string(payload), ":", 2)
	if len(creds) != 2 {
		ctx.Response().Header().Set("WWW-Authenticate",
			fmt.Sprintf("%srealm=\"Restricted\"", basicAuthPrefix))
		return "", "", ErrAuthData
	}
	return creds[0], creds[1], nil
}

func (a authenticator) hashPassword(pwd string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash
}

func (a authenticator) isPasswordsEqual(password []byte, hashPwd []byte) error {
	return bcrypt.CompareHashAndPassword(hashPwd, password)
}

func (h *HttpController) handlerSignUp(ctx echo.Context) error {
	userCfg := repo.UserConfig{}
	if err := ctx.Bind(&userCfg); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrBindingScheme.Error()})
	}

	if err := validateUser(userCfg); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}

	if _, err := h.store.GetUserByPassport(userCfg.Passport); err == nil {
		return ctx.JSON(http.StatusConflict,
			errorResponse{fmt.Sprintf("error of sign-up: %s", ErrUserExists)})
	}

	_, pwd, err := h.auth.extractCreds(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}
	userCfg.PwdHash = h.auth.hashPassword(pwd)

	return ctx.JSON(http.StatusCreated, h.store.CreateUser(userCfg))
}

func (h *HttpController) handlerLogin(ctx echo.Context) error {
	type authInfo struct {
		AccessToken string `json:"token"`
	}

	login, pwd, err := h.auth.extractCreds(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}

	user, err := h.store.GetUserByPassport(login)
	if err != nil {
		return ctx.JSON(http.StatusConflict,
			errorResponse{ErrInvalidLoginOrPwd.Error()})
	}

	if err := h.auth.isPasswordsEqual([]byte(pwd), user.PwdHash); err != nil {
		return ctx.JSON(http.StatusConflict,
			errorResponse{ErrInvalidLoginOrPwd.Error()})
	}

	token, err := h.auth.userToken.NewUserToken(user)
	if err != nil {
		h.logger.Error(err.Error())
		return ctx.JSON(http.StatusInternalServerError,
			errorResponse{ErrAuthFailed.Error()})
	}
	return ctx.JSON(http.StatusOK, authInfo{token})
}

func (h *HttpController) handlerVerifyToken(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := h.auth.extractRawToken(ctx)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{
				ErrAuthData.Error(),
			})
		}

		id, err := h.auth.userToken.VerifyUserJWT(token)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}

		_, err = h.store.GetUser(id)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, errorResponse{
				ErrInvalidToken.Error(),
			})
		}

		return handler(ctx)
	}
}

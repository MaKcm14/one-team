package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/mw/auth/token"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type httpError struct {
	code int
	resp errorResponse
}

func (a Authenticator) HandlerLogin(eCtx echo.Context) error {
	creds, httpErr := parseRequestForCreds(eCtx)
	if httpErr != nil {
		return eCtx.JSON(httpErr.code, httpErr.resp)
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	dto, err := a.authService.Login(ctx, creds)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrWrongPassword) || errors.Is(err, user.ErrRoleNotAssign) {
			return eCtx.JSON(http.StatusUnauthorized, errorResponse{
				Error: ErrInvalidAuthInfo.Error(),
			})
		}
		return eCtx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}
	return a.issueTokens(eCtx, dto)
}

func (a Authenticator) HandlerLogout(ctx echo.Context) error {
	session, err := a.session.Writer.Get(ctx.Request(), sessionIDCookieKey)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			ErrHandleRequest.Error(),
		})
	}
	rawSessionID := session.Values[sessionIDCookieKey]

	delete(session.Values, sessionIDCookieKey)

	sessionID, ok := rawSessionID.(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestInfo.Error(),
		})
	}
	a.tokens.RefreshTokens.Delete(sessionID)
	a.tokens.AccessTokens.Delete(sessionID)

	err = session.Save(ctx.Request(), ctx.Response().Writer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			ErrHandleRequest.Error(),
		})
	}
	return ctx.NoContent(http.StatusOK)
}

func parseRequestForCreds(ctx echo.Context) (user.Credentials, *httpError) {
	const basicAuth = "Basic "

	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, basicAuth) {
		return user.Credentials{}, &httpError{
			code: http.StatusUnauthorized,
			resp: errorResponse{
				ErrInvalidAuthHeader.Error(),
			},
		}
	}

	rawCreds, err := base64.StdEncoding.DecodeString(authHeader[len(basicAuth):])
	if err != nil {
		return user.Credentials{}, &httpError{
			code: http.StatusBadRequest,
			resp: errorResponse{
				ErrInvalidAuthHeader.Error(),
			},
		}
	}

	creds := strings.Split(string(rawCreds), ":")
	if len(creds) != 2 {
		return user.Credentials{}, &httpError{
			code: http.StatusBadRequest,
			resp: errorResponse{
				ErrWrongAuthInfo.Error(),
			},
		}
	}
	return user.Credentials{
		Login:    creds[0],
		Password: creds[1],
	}, nil
}

func (a Authenticator) issueTokens(ctx echo.Context, dto user.UserDTO) error {
	type tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	accessToken, err := a.acToken.IssueAccessToken(token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    token.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(token.AccessTokenTTL)),
		},
		UserData: user.Claims{
			Login: dto.User.Login,
			Role:  dto.Role,
		},
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}

	id, err := a.createSession()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}

	session, err := a.session.Writer.Get(ctx.Request(), sessionIDCookieKey)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}
	session.Values[sessionIDCookieKey] = id

	err = a.session.Writer.Save(ctx.Request(), ctx.Response().Writer, session)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{
			Error: ErrHandleRequest.Error(),
		})
	}

	refreshToken := token.IssueRefreshToken(128)
	a.tokens.AccessTokens.Set(id, accessToken, 0)
	a.tokens.RefreshTokens.Set(id, refreshToken, 0)

	return ctx.JSON(http.StatusOK, tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

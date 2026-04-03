package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type httpError struct {
	code int
	resp server.ErrorResponse
}

func (a Authenticator) HandlerLogin(eCtx echo.Context) error {
	creds, httpErr := parseRequestForCreds(eCtx)
	if httpErr != nil {
		a.log.Warn(fmt.Sprintf("Warn of parsing the creds"))
		return eCtx.JSON(httpErr.code, httpErr.resp)
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	dto, err := a.authService.Login(ctx, creds)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrWrongPassword) || errors.Is(err, user.ErrRoleNotAssign) {
			a.log.Error(fmt.Sprintf("Error of authentication: %s", err))
			return eCtx.JSON(http.StatusUnauthorized, server.ErrorResponse{
				Error: ErrInvalidAuthInfo.Error(),
			})
		}
		a.log.Error(fmt.Sprintf("Error of app-module: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}

	sid, err := a.createSession()
	if err != nil {
		a.log.Error(fmt.Sprintf("Error of creating the session: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return a.issueTokens(eCtx, sid, user.UserSession{
		UserClaims: user.Claims{
			Login: dto.User.Login,
			Role:  dto.Role,
		},
	})
}

func (a Authenticator) HandlerLogout(ctx echo.Context) error {
	claims, err := ExtractClaimsFromCtx(ctx)
	if err != nil {
		a.log.Warn(fmt.Sprintf("Warn of extracting the claims: %s", err))
		return ctx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	a.tokens.AccessTokens.Delete(claims.SessionID)
	a.tokens.RefreshTokens.Delete(claims.SessionID)
	a.session.Sessions.Delete(claims.SessionID)

	a.session.Sessions.Delete(claims.UserData.Login)

	return ctx.NoContent(http.StatusOK)
}

func (a Authenticator) HandlerRefresh(ctx echo.Context) error {
	type tokensRequest struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	tokens := tokensRequest{}
	if err := ctx.Bind(&tokens); err != nil {
		a.log.Error(fmt.Sprintf("Error of binding the body-request at refresh-operation: %s", err))
		return ctx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	claims, err := ExtractClaimsFromCtx(ctx)
	if err != nil {
		a.log.Warn(fmt.Sprintf("Warn of extracting the claims: %s", err))
		return ctx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	hashRefreshToken, err := a.tokens.GetHashRefreshToken(claims.SessionID)
	if err != nil {
		a.log.Warn(fmt.Sprintf("Warn of refresh-token storage: %s", err))
		return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
			Error: ErrRefreshTokenNotValid.Error(),
		})
	}

	err = a.refToken.CheckRefreshToken(hashRefreshToken, tokens.RefreshToken)
	if err != nil {
		a.log.Warn(fmt.Sprintf("Warn of refresh-token: it's not valid: %s", err))
		return ctx.JSON(http.StatusUnauthorized, server.ErrorResponse{
			Error: ErrRefreshTokenNotValid.Error(),
		})
	}

	a.tokens.RefreshTokens.Delete(claims.SessionID)

	userSession, err := a.session.GetSession(claims.SessionID)
	if err != nil {
		a.log.Warn(fmt.Sprintf("Warn of getting the session: %s", err))
		return ctx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return a.issueTokens(ctx, claims.SessionID, userSession)
}

func (a Authenticator) HandlerSignUp(eCtx echo.Context) error {
	dto := user.UserSignUpDTO{}
	if err := eCtx.Bind(&dto); err != nil {
		a.log.Error(
			fmt.Sprintf("Error of sign-up: %s", err),
		)
		return eCtx.JSON(http.StatusBadRequest, server.ErrorResponse{
			Error: server.ErrRequestInfo.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(eCtx.Request().Context(), 8*time.Second)
	defer cancel()

	err := a.authService.SignUp(ctx, dto)
	if err != nil {
		if errors.Is(err, user.ErrSignUp) {
			a.log.Warn(fmt.Sprintf("Warn of sign-up: %s", err))
			if errors.Is(err, user.ErrUserAlreadyExists) {
				return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
					Error: ErrSignUpUserExists.Error(),
				})
			} else if errors.Is(err, user.ErrVerifyPassword) {
				if errors.Is(err, user.ErrPasswordLength) {
					return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
						Error: "password length must be at least 9 symbols and less than 17",
					})
				} else if errors.Is(err, user.ErrPasswordSymbols) {
					return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
						Error: "password must contain at least 2 symbols from the list '@, #, _, !, $, ?'",
					})
				}
			} else if errors.Is(err, user.ErrRoleNotFound) {
				return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
					Error: "the set user's role doesn't exist",
				})
			}
		}
		a.log.Error(fmt.Sprintf("Error of sign-up: %s", err))
		return eCtx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}
	return eCtx.NoContent(http.StatusCreated)
}

func parseRequestForCreds(ctx echo.Context) (user.Credentials, *httpError) {
	const basicAuth = "Basic "

	authHeader := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, basicAuth) {
		return user.Credentials{}, &httpError{
			code: http.StatusUnauthorized,
			resp: server.ErrorResponse{
				Error: ErrInvalidAuthHeader.Error(),
			},
		}
	}

	rawCreds, err := base64.StdEncoding.DecodeString(authHeader[len(basicAuth):])
	if err != nil {
		return user.Credentials{}, &httpError{
			code: http.StatusBadRequest,
			resp: server.ErrorResponse{
				Error: ErrInvalidAuthHeader.Error(),
			},
		}
	}

	creds := strings.Split(string(rawCreds), ":")
	if len(creds) != 2 {
		return user.Credentials{}, &httpError{
			code: http.StatusBadRequest,
			resp: server.ErrorResponse{
				Error: ErrWrongAuthInfo.Error(),
			},
		}
	}
	return user.Credentials{
		Login:    creds[0],
		Password: creds[1],
	}, nil
}

// issuTokens defines the logic of issuing the tokens for the given session.
func (a Authenticator) issueTokens(ctx echo.Context, sid string, userSession user.UserSession) error {
	type tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	accessTokenID, _ := uuid.NewRandom()
	accessToken, err := a.acToken.IssueAccessToken(token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    token.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(token.AccessTokenTTL)),
			ID:        accessTokenID.String(),
		},
		UserData:  userSession.UserClaims,
		SessionID: sid,
	})
	if err != nil {
		a.log.Error(fmt.Sprintf("Error of issue the token: %s", err))
		return ctx.JSON(http.StatusInternalServerError, server.ErrorResponse{
			Error: server.ErrHandleRequest.Error(),
		})
	}

	refreshToken := a.refToken.IssueRefreshToken(64)
	refreshTokenHash, err := a.refToken.HashRefreshToken(refreshToken)

	a.tokens.AccessTokens.Set(sid, accessTokenID.String(), 0)
	a.tokens.RefreshTokens.Set(sid, string(refreshTokenHash), 0)
	a.session.Sessions.Set(sid, userSession, 0)
	a.session.SetSIDForLogin(userSession.UserClaims.Login, sid, 0)

	return ctx.JSON(http.StatusOK, tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

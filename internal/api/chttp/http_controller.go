package chttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo"

	"auth-train/test/internal/api/chttp/auth"
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

type errorResponse struct {
	Error string `json:"error"`
}

type HttpController struct {
	e     *echo.Echo
	auth  auth.Authenticator
	store repo.Repository

	logger *slog.Logger
}

func New(logger *slog.Logger, authConf tokens.AuthJWTConfig) HttpController {
	return HttpController{
		e:      echo.New(),
		auth:   auth.NewAuthenticator(authConf),
		store:  repo.NewRepository(logger),
		logger: logger,
	}
}

func (h *HttpController) Run(socket string) error {
	h.configPath()
	if err := h.e.Start(socket); err != nil {
		errRet := fmt.Errorf("%s: %s", ErrStartController, err)
		h.logger.Error(err.Error())
		return errRet
	}
	return nil
}

func (h *HttpController) configPath() {
	h.e.File("/index.html", "../../web/api/static/index.html")

	h.e.GET("/get/user", h.handlerGetUser, h.handlerVerifyToken)
	h.e.GET("/get/user/list", h.handlerGetUserList, h.handlerVerifyToken)

	h.e.POST("/signup", h.handlerSignUp)
	h.e.POST("/login", h.handlerLogin)

	h.e.DELETE("/delete/user", h.handlerDeleteUser, h.handlerVerifyToken)

	h.e.PATCH("/set/money", h.handlerSetMoney, h.handlerVerifyToken)
}

func (h *HttpController) handlerDeleteUser(ctx echo.Context) error {
	id, err := validateUserID(ctx.QueryParam(userIDParamName))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestQueryParam.Error(),
		})
	}
	h.store.DeleteUser(id)
	return ctx.NoContent(http.StatusResetContent)
}

func (h *HttpController) handlerGetUser(ctx echo.Context) error {
	id, err := validateUserID(ctx.QueryParam(userIDParamName))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestQueryParam.Error(),
		})
	}

	user, err := h.store.GetUser(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errorResponse{
			ErrResourceNotFound.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h *HttpController) handlerGetUserList(ctx echo.Context) error {
	type usersResponse struct {
		Users []entity.User `json:"users"`
	}
	return ctx.JSON(http.StatusOK, usersResponse{h.store.GetUsers()})
}

func (h *HttpController) handlerSetMoney(ctx echo.Context) error {
	type request struct {
		ID    entity.UserID `json:"user_id"`
		Money float64       `json:"money"`
	}
	var req request

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{
			ErrRequestBody.Error(),
		})
	}

	user, err := h.store.SetMoney(req.ID, req.Money)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errorResponse{
			ErrResourceNotFound.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, user)
}

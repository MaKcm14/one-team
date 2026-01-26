package chttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo"

	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

type errorResponse struct {
	Error string `json:"error"`
}

type HttpController struct {
	e     *echo.Echo
	store repo.BankRepository

	logger *slog.Logger
}

func New(logger *slog.Logger) HttpController {
	return HttpController{
		e:      echo.New(),
		store:  repo.NewBankRepository(logger),
		logger: logger,
	}
}

func (h HttpController) Run(socket string) error {
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

	h.e.GET("/get/user", h.handlerGetUser)
	h.e.GET("/get/user/list", h.handlerGetUserList)

	h.e.POST("/create/user", h.handlerCreateUser)

	h.e.DELETE("/delete/user", h.handlerDeleteUser)

	h.e.PATCH("/set/money", h.handlerSetMoney)
}

func (h *HttpController) handlerCreateUser(ctx echo.Context) error {
	userCfg := repo.UserConfig{}
	if err := ctx.Bind(&userCfg); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrBindingScheme.Error()})
	}

	if err := validateUser(userCfg); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}
	return ctx.JSON(http.StatusCreated, h.store.CreateUser(userCfg))
}

func (h *HttpController) handlerDeleteUser(ctx echo.Context) error {
	id, err := validateUserID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}
	h.store.DeleteUser(id)
	return ctx.JSON(http.StatusResetContent, nil)
}

func (h *HttpController) handlerGetUser(ctx echo.Context) error {
	id, err := validateUserID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
	}

	user, err := h.store.GetUser(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errorResponse{ErrRequestData.Error()})
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
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrBindingScheme.Error()})
	}

	user, err := h.store.SetMoney(req.ID, req.Money)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, errorResponse{ErrRequestData.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

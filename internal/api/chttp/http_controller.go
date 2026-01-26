package chttp

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo"
)

type HttpController struct {
	e      *echo.Echo
	logger *slog.Logger
}

func New(logger *slog.Logger) HttpController {
	return HttpController{
		e:      echo.New(),
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

func (h HttpController) configPath() {

}

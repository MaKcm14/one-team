package chttp

import (
	"fmt"
	"log/slog"

	"github.com/MaKcm14/one-team/internal/api"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	e   *echo.Echo
	log *slog.Logger

	cfg config.ControllerConfig
}

func New(log *slog.Logger, cfg config.ControllerConfig) *Controller {
	return &Controller{
		e:   echo.New(),
		log: log,
		cfg: cfg,
	}
}

func (c Controller) Run() error {
	c.log.Info("starting the app-controller")

	c.configEndpoints()

	if err := c.e.Start(c.cfg.Socket); err != nil {
		return fmt.Errorf("%w: %s", api.ErrStartController, err)
	}
	return nil
}

func (c Controller) configEndpoints() {

}

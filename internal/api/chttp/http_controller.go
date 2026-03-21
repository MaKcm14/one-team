package chttp

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/mw"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Controller struct {
	e   *echo.Echo
	log *slog.Logger

	cfg    config.ControllerConfig
	authMW auth.Authenticator
}

func New(
	log *slog.Logger,
	cfg config.ControllerConfig,
	authService user.IAuthService,
) *Controller {
	return &Controller{
		e:      echo.New(),
		log:    log,
		cfg:    cfg,
		authMW: auth.NewMW(log, cfg.AuthCfg, authService),
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
	c.e.Use(
		mw.Recovery(c.log),
		mw.LoggerMW(c.log),
		c.authMW.DebugPrintCaches(),
	)

	adminGroup := c.e.Group("/admin") //c.authMW.VerifyAccessTokenMW())
	{
		adminGroup.POST("/signup", c.authMW.HandlerSignUp)
	}

	clientGroup := c.e.Group("/client", c.authMW.VerifyAccessTokenMW())
	{
		clientGroup.POST("/logout", c.authMW.HandlerLogout)
	}

	authGroup := c.e.Group("/auth")
	{
		authGroup.POST("/login", c.authMW.HandlerLogin)
		authGroup.POST("/token/refresh", c.authMW.HandlerRefresh)
	}
}

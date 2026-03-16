package app

import (
	"fmt"

	"github.com/MaKcm14/one-team/internal/api/chttp"
	"github.com/MaKcm14/one-team/internal/app/logger"
	"github.com/MaKcm14/one-team/internal/config"
)

type App struct {
	log logger.Logger

	contr *chttp.Controller
}

func New() App {
	cfg, err := config.New(
		config.DefaultConfigFilePath,
	)
	if err != nil {
		panic(err)
	}

	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	return App{
		log:   log,
		contr: chttp.New(log.Instance(), cfg.ControllerCfg),
	}
}

func (a *App) Run() {
	defer a.log.Close()
	defer a.log.Info("STOP THE APP")

	a.log.Info("STARTING THE APP")
	if err := a.contr.Run(); err != nil {
		a.log.Error(fmt.Sprintf("error of starting the controller: %s", err))
	}
}

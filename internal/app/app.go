package app

import (
	"fmt"

	"github.com/MaKcm14/one-team/internal/api/chttp"
	"github.com/MaKcm14/one-team/internal/app/logger"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/repository/persistent/postgres"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
	"github.com/MaKcm14/one-team/internal/services/usecase/root"
	"github.com/MaKcm14/one-team/internal/services/usecase/user/auth"
)

type App struct {
	log logger.Logger

	contr *chttp.Controller
	repo  postgres.Repository
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

	repo, err := postgres.NewRepository(log, cfg.DBCfg)
	if err != nil {
		panic(err)
	}

	return App{
		log: log,
		contr: chttp.New(
			log.Instance(),
			cfg.ControllerCfg,
			auth.NewInteractor(log, cfg.ControllerCfg.AuthCfg, repo),
			root.NewInteractor(repo),
			employee.NewInteractor(repo),
			division.NewInteractor(repo),
		),
		repo: repo,
	}
}

func (a *App) Run() {
	defer a.log.Close()
	defer a.repo.Close()
	defer a.log.Info("STOP THE APP")

	a.log.Info("STARTING THE APP")
	if err := a.contr.Run(); err != nil {
		a.log.Error(fmt.Sprintf("error of starting the controller: %s", err))
	}
}

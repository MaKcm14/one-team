package app

import (
	"log/slog"
	"os"

	"auth-train/test/internal/api/chttp"
	"auth-train/test/internal/config"
)

type AuthService struct {
	logger     *slog.Logger
	conf       config.AuthServiceConfig
	controller chttp.HttpController
}

func NewAuthService() AuthService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	conf := config.NewAuthServiceConfig(
		logger,
		config.ConfigSocket,
	)
	contr := chttp.New(logger)

	return AuthService{
		logger:     logger,
		conf:       conf,
		controller: contr,
	}
}

func (a AuthService) Start() {
	defer a.logger.Info("AUTH_SERVICE_STOP")

	a.logger.Info("AUTH_SERVICE_START")
	if err := a.controller.Run(a.conf.Socket); err != nil {
		panic(err)
	}
}

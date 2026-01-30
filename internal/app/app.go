package app

import (
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"auth-train/test/internal/api/chttp"
	"auth-train/test/internal/api/chttp/auth/tokens"
	"auth-train/test/internal/config"
)

type AuthService struct {
	conf       config.AuthServiceConfig
	controller chttp.HttpController

	logger *slog.Logger
}

func NewAuthService() AuthService {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	conf := config.MustAuthServiceConfig(
		logger,
		config.WithSocket(),
		config.WithSecret(),
	)
	contr := chttp.New(
		logger,
		tokens.AuthJWTConfig{
			Secret: []byte(conf.Secret),
			Method: jwt.SigningMethodHS256,
		},
	)

	return AuthService{
		conf:       conf,
		controller: contr,
		logger:     logger,
	}
}

func (a AuthService) Start() {
	defer a.logger.Info("AUTH_SERVICE_STOP")

	a.logger.Info("AUTH_SERVICE_START")
	if err := a.controller.Run(a.conf.Socket); err != nil {
		panic(err)
	}
}

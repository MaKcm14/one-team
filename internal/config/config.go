package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)

type AuthServiceConfig struct {
	Socket string
	Secret string
}

func NewAuthServiceConfig(logger *slog.Logger, opts ...AuthConfigOpt) AuthServiceConfig {
	logger.Info("AUTH_SERVICE_CONFIGURING_START")

	conf := AuthServiceConfig{}
	godotenv.Load("../../.env")

	for _, opt := range opts {
		if err := opt(&conf); err != nil {
			errRet := fmt.Errorf("%s: %s", ErrAuthServiceConfiguration, err)
			logger.Error(err.Error())
			panic(errRet.Error())
		}
	}
	return conf
}

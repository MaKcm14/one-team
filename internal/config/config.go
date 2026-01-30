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

func MustAuthServiceConfig(logger *slog.Logger, opts ...AuthConfigOpt) AuthServiceConfig {
	logger.Info("AUTH_SERVICE_CONFIGURING_START")

	conf := AuthServiceConfig{}
	godotenv.Load("../../.env")

	for _, opt := range opts {
		if err := opt(&conf); err != nil {
			errRet := fmt.Errorf("%s: %s", ErrServiceConfig, err)
			logger.Error(err.Error())
			panic(errRet.Error())
		}
	}
	return conf
}

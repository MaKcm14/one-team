package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultConfigFilePath = "./config.yaml"

type ControllerConfig struct {
	Socket     string `yaml:"socket" env-required:"true"`
	GlobalSalt int    `yaml:"salt" env-required:"true"`
}

type Config struct {
	ControllerCfg ControllerConfig `yaml:"controller"`
}

func New(path string) (Config, error) {
	var conf Config
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		return Config{}, fmt.Errorf("%w: %s", ErrConfigParse, err)
	}
	return conf, nil
}

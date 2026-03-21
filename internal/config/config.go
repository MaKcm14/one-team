package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const DefaultConfigFilePath = "./config.yaml"

type DBConfig struct {
	DSN string `yaml:"dsn" env-required:"true"`
}

type AuthConfig struct {
	Secret        string `yaml:"secret" env-required:"true"`
	SessionKey    string `yaml:"session_key" env-required:"true"`
	GlobalPwdSalt int    `yaml:"pwd_salt" env-required:"true"`
	TokenSalt     int    `yaml:"token_salt" env-required:"true"`
}

type ControllerConfig struct {
	Socket  string     `yaml:"socket" env-required:"true"`
	AuthCfg AuthConfig `yaml:"auth" env-required:"true"`
}

type Config struct {
	DBCfg         DBConfig         `yaml:"db"`
	ControllerCfg ControllerConfig `yaml:"controller"`
}

func New(path string) (Config, error) {
	var conf Config
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		return Config{}, fmt.Errorf("%w: %s", ErrConfigParse, err)
	}
	return conf, nil
}

package config

import (
	"fmt"
	"os"
	"regexp"
)

const (
	SOCKET_SETTING_NAME = "SOCKET"
	SECRET_SETTING_NAME = "SECRET"
)

type AuthConfigOpt func(*AuthServiceConfig) error

func ConfigSocket(conf *AuthServiceConfig) error {
	socket := os.Getenv(SOCKET_SETTING_NAME)
	socketReg := regexp.MustCompile(
		`^(localhost|\d+\.\d+\.\d+\.\d+):\d+$`,
	)

	if len(socket) == 0 {
		return fmt.Errorf("error of %s: it can't be empty", SOCKET_SETTING_NAME)
	} else if !socketReg.MatchString(socket) {
		return fmt.Errorf("error of %s: it doesn't match the correct form of IP:PORT",
			SOCKET_SETTING_NAME)
	}
	conf.Socket = socket
	return nil
}

func ConfigSecret(conf *AuthServiceConfig) error {
	secret := os.Getenv(SECRET_SETTING_NAME)

	if len(secret) < 256 {
		return fmt.Errorf("error of %s: can't be less than 256bits", SECRET_SETTING_NAME)
	}
	conf.Secret = secret

	return nil
}

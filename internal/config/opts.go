package config

import (
	"fmt"
	"os"
	"regexp"
)

const (
	SOCKET_SETTING_NAME = "SOCKET"
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

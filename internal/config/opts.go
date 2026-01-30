package config

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"sync"
)

const (
	SOCKET_SETTING_NAME = "SOCKET"
	SECRET_SETTING_NAME = "SECRET"
)

type AuthConfigOpt func(*AuthServiceConfig) error

func WithSocket() AuthConfigOpt {
	return func(conf *AuthServiceConfig) error {
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
}

func genSecret() string {
	const size = 256
	var (
		buff     = make([]byte, size)
		wg       = sync.WaitGroup{}
		poolSize = runtime.GOMAXPROCS(0)
		intNum   = size / poolSize
	)

	wg.Add(poolSize)
	for i := 0; i != poolSize; i++ {
		part := buff[i*intNum:]
		if i != poolSize-1 {
			part = buff[i*intNum : (i+1)*intNum]
		}

		go func() {
			defer wg.Done()
			genSecretPart(part)
		}()
	}
	wg.Wait()

	return string(buff)
}

func genSecretPart(slice []byte) {
	for i := 0; i != len(slice); i++ {
		val := []rune(fmt.Sprintf("%X", rand.Intn(16)))[0]
		slice[i] = byte(val)
	}
}

func WithSecret() AuthConfigOpt {
	return func(conf *AuthServiceConfig) error {
		secret := os.Getenv(SECRET_SETTING_NAME)

		if len(secret) < 256 && secret != "auto" {
			return fmt.Errorf("error of %s: can't be less than 256bits", SECRET_SETTING_NAME)
		} else if secret == "auto" {
			secret = genSecret()
		}
		conf.Secret = secret

		return nil
	}
}

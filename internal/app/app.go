package app

import (
	"fmt"

	"github.com/MaKcm14/one-team/internal/config"
)

type App struct {
}

func New() App {
	conf, err := config.New(
		config.DefaultConfigFilePath,
	)
	if err != nil {
		panic(err)
	}
	_ = conf

	fmt.Println(conf)

	return App{}
}

func (a *App) Run() {

}

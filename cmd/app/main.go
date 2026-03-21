package main

import "github.com/MaKcm14/one-team/internal/app"

func main() {
	service := app.New()
	service.Run()
}

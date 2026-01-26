package main

import "auth-train/test/internal/app"

func main() {
	service := app.NewAuthService()
	service.Start()
}

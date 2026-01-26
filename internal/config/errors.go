package config

import "errors"

var (
	ErrAuthServiceConfiguration = errors.New(
		"config: error of configuration the auth-service",
	)
)

package repo

import "errors"

var (
	ErrUserNotExist = errors.New("repo: error of op with the unexisting user")
)

package user

import "errors"

var (
	ErrWrongPassword = errors.New("user: error of comparing the passwords: not equal")
	ErrUserNotFound  = errors.New("user: error of searching the user: not found")
	ErrRepoInteract  = errors.New("user: error of interact with the repository")
)

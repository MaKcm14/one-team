package user

import "errors"

var (
	ErrWrongPassword = errors.New("user: error of comparing the passwords: not equal")
	ErrUserNotFound  = errors.New("user: error of searching the user: not found")
	ErrRepoInteract  = errors.New("user: error of interact with the repository")
	ErrRoleNotAssign = errors.New("user: error of getting the user's role: not assign")
)

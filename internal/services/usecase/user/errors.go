package user

import "errors"

var (
	ErrHashPassword  = errors.New("user: error of hashing the password")
	ErrWrongPassword = errors.New("user: error of comparing the passwords: not equal")
	ErrUserNotFound  = errors.New("user: error of searching the user: not found")
	ErrRepoInteract  = errors.New("user: error of interact with the repository")
	ErrRoleNotAssign = errors.New("user: error of getting the user's role: not assign")
	ErrRoleNotFound  = errors.New("user: error of searching the role: not found")
	ErrSignUp        = errors.New("user: error of sign-up the user")
)

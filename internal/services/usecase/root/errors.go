package root

import "errors"

var (
	ErrRepoInteract        = errors.New("root: error of interaction with repository")
	ErrUnableToDeleteAdmin = errors.New("root: error of deleting the admin: restict")
	ErrUserNotFound        = errors.New("root: error of searching the user: not found")
)

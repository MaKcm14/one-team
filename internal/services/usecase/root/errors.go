package root

import "errors"

var (
	ErrRepoInteract            = errors.New("root: error of interaction with repository")
	ErrRoleNotFound            = errors.New("root: error of searching the role: not found")
	ErrUnableToDeleteAdmin     = errors.New("root: error of deleting the admin: restict")
	ErrUnableToChangeAdminRole = errors.New("root: error of changing the admin role: restrict")
	ErrUserNotFound            = errors.New("root: error of searching the user: not found")
)

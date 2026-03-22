package persistent

import "errors"

var (
	ErrCitizenshipNotFound = errors.New("persistent: error of searching the citizenship: not found")
	ErrConnWithDB          = errors.New("persistent: error of setting the connection with the DB")
	ErrQueryExec           = errors.New("persistent: error of executing the query")
	ErrUserNotFound        = errors.New("persistent: error of searching the user: not found")
	ErrRoleNotAssign       = errors.New("persistent: error of getting the role: not assign")
	ErrRoleNotFound        = errors.New("persistent: error of getting the role: not found")
	ErrTitleNotFound       = errors.New("persistent: error of searching the title: not found")
)

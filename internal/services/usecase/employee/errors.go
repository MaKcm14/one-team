package employee

import "errors"

var (
	ErrEmployeeExists      = errors.New("employee: error of creating employee: already exists")
	ErrCitizenshipNotFound = errors.New("employee: error of searching the citizenship: not found")
	ErrRepoInteract        = errors.New("employee: error of interaction with repo")
	ErrTitleNotFound       = errors.New("employee: error of searching the title: not found")
)

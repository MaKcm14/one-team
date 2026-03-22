package employee

import "errors"

var (
	ErrCitizenshipNotFound = errors.New("employee: error of searching the citizenship: not found")
	ErrRepoInteract        = errors.New("employee: error of interaction with repo")
	ErrTitleNotFound       = errors.New("employee: error of searching the title: not found")
)

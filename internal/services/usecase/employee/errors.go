package employee

import "errors"

var (
	ErrEmployeeExists      = errors.New("employee: error of creating employee: already exists")
	ErrEmployeeNotFound    = errors.New("employee: error of searching the employee: not found")
	ErrCitizenshipNotFound = errors.New("employee: error of searching the citizenship: not found")
	ErrRepoInteract        = errors.New("employee: error of interaction with repo")
	ErrReportCreating      = errors.New("employee: error of creating a report")
	ErrTitleNotFound       = errors.New("employee: error of searching the title: not found")
)

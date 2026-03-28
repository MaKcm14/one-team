package division

import "errors"

var (
	ErrDivisionExists         = errors.New("division: error of creating the division: it already exists")
	ErrDivisionNotFound       = errors.New("division: error of searching the division: not found")
	ErrRepoInteract           = errors.New("division: error of interact with repo")
	ErrWrongDivisionsRelation = errors.New("division: error of inter-divisions types' relation: it's wrong")
)

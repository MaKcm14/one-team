package division

import "errors"

var (
	ErrDivisionExists          = errors.New("division: error of creating the division: it already exists")
	ErrDivisionNotEmpty        = errors.New("division: the division is not empty")
	ErrDivisionNotFound        = errors.New("division: error of searching the division: not found")
	ErrDivisionIsSuperdivision = errors.New("division: the division is superdivision")
	ErrRepoInteract            = errors.New("division: error of interact with repo")
	ErrWrongDivisionsRelation  = errors.New("division: error of inter-divisions types' relation: it's wrong")
)

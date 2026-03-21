package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type Interactor struct {
}

func NewInteractor() Interactor {
	return Interactor{}
}

func (e Interactor) CreateEmployee(ctx context.Context, employee entity.Employee) error {
	return nil
}

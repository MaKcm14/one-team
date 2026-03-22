package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type Interactor struct {
	divisionRepo IDivisionRepo
}

func NewInteractor(divisionRepo IDivisionRepo) Interactor {
	return Interactor{
		divisionRepo: divisionRepo,
	}
}

func (d Interactor) GetDivisions(ctx context.Context) ([]entity.Division, error) {
	return nil, nil
}

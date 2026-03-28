package division

import (
	"context"
	"fmt"

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

func (d Interactor) GetDivisions(ctx context.Context, filter Filters) ([]entity.Division, error) {
	divisions, err := d.divisionRepo.GetDivisionsByName(ctx, filter.Names)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return divisions, nil
}

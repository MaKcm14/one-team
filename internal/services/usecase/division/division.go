package division

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
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

func (d Interactor) CreateDivision(ctx context.Context, div entity.Division) error {
	err := d.divisionRepo.IsDivisionExists(ctx, div)
	if err == nil {
		return ErrDivisionExists
	} else if !errors.Is(err, persistent.ErrDivisionNotFound) {
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	if div.Type == entity.DivisionTypeName {
		err = d.divisionRepo.CreateDivisionOfDivisionType(ctx, div)
	} else {
		var supDiv entity.Division

		supDiv, err = d.divisionRepo.GetDivisionByID(ctx, div.SuperdivisionID)
		if err != nil {
			if errors.Is(err, persistent.ErrDivisionNotFound) {
				return fmt.Errorf("superdivision search process: %w: %s", ErrDivisionNotFound, err)
			}
			return fmt.Errorf("%w: %s", ErrRepoInteract, err)
		}

		if !entity.IsDivisionTypeRelationCorrect(div.Type, supDiv.Type) {
			return ErrWrongDivisionsRelation
		}
		err = d.divisionRepo.CreateDivisionOfNotDivisionType(ctx, div)
	}

	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return nil
}

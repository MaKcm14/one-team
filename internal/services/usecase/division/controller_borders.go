package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionServiceReader interface {
	GetDivisions(ctx context.Context, filters Filters) ([]entity.Division, error)
}

type IDivisionServiceWriter interface {
	CreateDivision(ctx context.Context, div entity.Division) error
	DeleteDivision(ctx context.Context, id int) error
}

type IDivisionService interface {
	IDivisionServiceReader
	IDivisionServiceWriter
}

package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionServiceReader interface {
	GetDivisions(ctx context.Context, filters Filters) ([]entity.Division, error)
	GetSalaryStatisticsOfDivision(ctx context.Context, id int) (SalaryStatistics, error)
}

type IDivisionServiceWriter interface {
	CreateDivision(ctx context.Context, div entity.Division) error
	DeleteDivision(ctx context.Context, id int) error
	UpdateDivision(ctx context.Context, div entity.Division) error
}

type IDivisionService interface {
	IDivisionServiceReader
	IDivisionServiceWriter
}

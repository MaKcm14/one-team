package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionRepoReader interface {
	CheckDivisionIsSuperdivision(ctx context.Context, id int) error
	GetDivisionsByName(ctx context.Context, filter NameFilter) ([]entity.Division, error)
	GetDivisionByID(ctx context.Context, id int) (entity.Division, error)
	IsDivisionExistsByName(ctx context.Context, div entity.Division) error
	IsDivisionExistsByID(ctx context.Context, id int) error
	IsDivisionEmpty(ctx context.Context, id int) error
	GetSalaryStatisticsOfDivision(ctx context.Context, id int) (SalaryStatistics, error)
	GetMinStateSizeDivisions(ctx context.Context, divType entity.DivisionType) ([]entity.Division, error)
	GetMaxStateSizeDivisions(ctx context.Context, divType entity.DivisionType) ([]entity.Division, error)
}

type IDivisionRepoWriter interface {
	CreateDivisionOfDivisionType(ctx context.Context, div entity.Division) error
	CreateDivisionOfNotDivisionType(ctx context.Context, div entity.Division) error
	UpdateDivisionOfDivisionType(ctx context.Context, div entity.Division) error
	UpdateDivisionOfNotDivisionType(ctx context.Context, div entity.Division) error
	DeleteDivisionByID(ctx context.Context, id int) error
}

type IDivisionRepo interface {
	IDivisionRepoReader
	IDivisionRepoWriter
}

package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionRepoReader interface {
	GetDivisionsByName(ctx context.Context, filter NameFilter) ([]entity.Division, error)
	GetDivisionByID(ctx context.Context, id int) (entity.Division, error)
	IsDivisionExists(ctx context.Context, div entity.Division) error
}

type IDivisionRepoWriter interface {
	CreateDivisionOfDivisionType(ctx context.Context, div entity.Division) error
	CreateDivisionOfNotDivisionType(ctx context.Context, div entity.Division) error
}

type IDivisionRepo interface {
	IDivisionRepoReader
	IDivisionRepoWriter
}

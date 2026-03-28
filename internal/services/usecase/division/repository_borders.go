package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionRepoReader interface {
	GetDivisionsByName(ctx context.Context, filter NameFilter) ([]entity.Division, error)
}

type IDivisionRepo interface {
	IDivisionRepoReader
}

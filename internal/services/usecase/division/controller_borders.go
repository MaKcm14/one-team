package division

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type IDivisionServiceReader interface {
	GetDivisions(ctx context.Context) ([]entity.Division, error)
}

type IDivisionService interface {
	IDivisionServiceReader
}

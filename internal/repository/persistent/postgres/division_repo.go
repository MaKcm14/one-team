package postgres

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
)

type divisionRepo struct {
	client *postgresClient
}

func (d divisionRepo) GetDivisions(ctx context.Context) ([]entity.Division, error) {
	return nil, nil
}

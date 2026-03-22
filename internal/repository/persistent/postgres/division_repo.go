package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
)

type divisionRepo struct {
	client *postgresClient
}

const getDivisionsQuery = `
SELECT id, name, type, state_size, superdivision_id
FROM usecase.divisions;
`

func (d divisionRepo) GetDivisions(ctx context.Context) ([]entity.Division, error) {
	res, err := d.client.conn.Query(ctx, getDivisionsQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	divisions := make([]entity.Division, 0, 30)
	for res.Next() {
		var division entity.Division
		err := res.Scan(
			&division.ID,
			&division.Name,
			&division.Type,
			&division.StateSize,
			&division.SuperdivisionID,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		divisions = append(divisions, division)
	}
	return divisions, nil
}

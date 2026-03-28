package postgres

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/division"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
)

type divisionRepo struct {
	client *postgresClient
}

const getDivisionsQuery = `
SELECT id, name, type, state_size, superdivision_id
FROM usecase.divisions
WHERE name LIKE $1 AND type LIKE $2;
`

func (d divisionRepo) GetDivisionsByName(ctx context.Context, filter division.NameFilter) ([]entity.Division, error) {
	res, err := d.client.conn.Query(
		ctx,
		getDivisionsQuery,
		as(filter.Name),
		as(string(filter.Type)),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	divisions := make([]entity.Division, 0, 30)
	for res.Next() {
		var (
			division        entity.Division
			superdivisionID sql.NullInt64
		)
		err := res.Scan(
			&division.ID,
			&division.Name,
			&division.Type,
			&division.StateSize,
			&superdivisionID,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		division.SuperdivisionID = int(superdivisionID.Int64)

		divisions = append(divisions, division)
	}
	return divisions, nil
}

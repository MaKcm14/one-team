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

const isDivisionExistsQuery = `
SELECT COUNT(*)
FROM usecase.divisions
WHERE name=$1 AND type=$2;
`

func (d divisionRepo) IsDivisionExists(ctx context.Context, div entity.Division) error {
	res, err := d.client.conn.Query(
		ctx,
		isDivisionExistsQuery,
		div.Name,
		div.Type,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if count == 0 {
		return persistent.ErrDivisionNotFound
	}
	return nil
}

const createDivisionNotDivisionTypeQuery = `
INSERT INTO usecase.divisions (name, type, state_size, superdivision_id)
VALUES ($1, $2, $3, $4);
`

func (d divisionRepo) CreateDivisionOfNotDivisionType(ctx context.Context, div entity.Division) error {
	_, err := d.client.conn.Exec(
		ctx,
		createDivisionNotDivisionTypeQuery,
		div.Name,
		div.Type,
		div.StateSize,
		div.SuperdivisionID,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return nil
}

const createDivisionDivisionTypeQuery = `
INSERT INTO usecase.divisions (name, type, state_size)
VALUES ($1, $2, $3);
`

func (d divisionRepo) CreateDivisionOfDivisionType(ctx context.Context, div entity.Division) error {
	_, err := d.client.conn.Exec(
		ctx,
		createDivisionDivisionTypeQuery,
		div.Name,
		div.Type,
		div.StateSize,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return nil
}

const getDivisionByIDQuery = `
SELECT id, name, type, state_size, superdivision_id
FROM usecase.divisions
WHERE id=$1;
`

func (d divisionRepo) GetDivisionByID(ctx context.Context, id int) (entity.Division, error) {
	res, err := d.client.conn.Query(
		ctx,
		getDivisionByIDQuery,
		id,
	)
	if err != nil {
		return entity.Division{}, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return entity.Division{}, persistent.ErrDivisionNotFound
	}

	var (
		div             entity.Division
		superdivisionID sql.NullInt64
	)
	err = res.Scan(
		&div.ID,
		&div.Name,
		&div.Type,
		&div.StateSize,
		&superdivisionID,
	)
	if err != nil {
		return entity.Division{}, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	div.SuperdivisionID = int(superdivisionID.Int64)

	return div, nil
}

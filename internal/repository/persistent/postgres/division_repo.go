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

const isDivisionExistsByNameQuery = `
SELECT COUNT(*)
FROM usecase.divisions
WHERE name=$1 AND type=$2;
`

func (d divisionRepo) IsDivisionExistsByName(ctx context.Context, div entity.Division) error {
	res, err := d.client.conn.Query(
		ctx,
		isDivisionExistsByNameQuery,
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

const deleteDivisionByIDQuery = `
DELETE FROM usecase.divisions
WHERE id=$1;
`

func (d divisionRepo) DeleteDivisionByID(ctx context.Context, id int) error {
	res, err := d.client.conn.Exec(ctx, deleteDivisionByIDQuery, id)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrDivisionNotFound
	}
	return nil
}

const isDivisionEmptyQuery = `
SELECT COUNT(*)
FROM usecase.employees
WHERE unit_id=$1;
`

func (d divisionRepo) IsDivisionEmpty(ctx context.Context, id int) error {
	res, err := d.client.conn.Query(ctx, isDivisionEmptyQuery, id)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return persistent.ErrQueryExec
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if count != 0 {
		return persistent.ErrDivisionNotEmpty
	}
	return nil
}

const checkDivisionIsSuperdivisionQuery = `
SELECT COUNT(*)
FROM usecase.divisions
WHERE superdivision_id=$1;
`

func (d divisionRepo) CheckDivisionIsSuperdivision(ctx context.Context, id int) error {
	res, err := d.client.conn.Query(ctx, checkDivisionIsSuperdivisionQuery, id)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return persistent.ErrQueryExec
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if count == 0 {
		return persistent.ErrDivisionNotSuperdivision
	}
	return nil
}

const updateDivisionOfNotDivisionTypeQuery = `
UPDATE usecase.divisions
SET name=$1,
	type=$2,
	state_size=$3,
	superdivision_id=$4
WHERE id=$5;
`

func (d divisionRepo) UpdateDivisionOfNotDivisionType(ctx context.Context, div entity.Division) error {
	res, err := d.client.conn.Exec(
		ctx,
		updateDivisionOfNotDivisionTypeQuery,
		div.Name,
		div.Type,
		div.StateSize,
		div.SuperdivisionID,
		div.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrDivisionNotFound
	}
	return nil
}

const updateDivisionOfDivisionTypeQuery = `
UPDATE usecase.divisions
SET name=$1,
	type=$2,
	state_size=$3,
	superdivision_id=NULL
WHERE id=$4;
`

func (d divisionRepo) UpdateDivisionOfDivisionType(ctx context.Context, div entity.Division) error {
	res, err := d.client.conn.Exec(
		ctx,
		updateDivisionOfDivisionTypeQuery,
		div.Name,
		div.Type,
		div.StateSize,
		div.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if res.RowsAffected() == 0 {
		return persistent.ErrDivisionNotFound
	}
	return nil
}

const isDivisionExistsByIDQuery = `
SELECT COUNT(*)
FROM usecase.divisions
WHERE id=$1;
`

func (d divisionRepo) IsDivisionExistsByID(ctx context.Context, id int) error {
	res, err := d.client.conn.Query(ctx, isDivisionExistsByIDQuery, id)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return persistent.ErrQueryExec
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

package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
)

type employeeRepo struct {
	client *postgresClient
}

const getCitizenshipIDByNameQuery = `
SELECT id
FROM usecase.citizenships
WHERE name=$1;
`

func (e employeeRepo) getCitizenshipIDByName(ctx context.Context, name entity.Citizenship) (int, error) {
	res, err := e.client.conn.Query(ctx, getCitizenshipIDByNameQuery, name)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return 0, persistent.ErrCitizenshipNotFound
	}

	var id int
	if err := res.Scan(&id); err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return id, nil
}

const getTitleIDByName = `
SELECT id
FROM usecase.titles
WHERE name=$1;
`

func (e employeeRepo) getTitleIDByName(ctx context.Context, name entity.Title) (int, error) {
	res, err := e.client.conn.Query(ctx, getTitleIDByName, name)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return 0, persistent.ErrTitleNotFound
	}

	var id int
	if err := res.Scan(&id); err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return id, nil
}

const createEmployeeQuery = `
INSERT INTO usecase.employees (tin_num, snils_num, passport_data, phone_num, first_name, last_name, patronymic, address, title_id, hiring_date, unit_id, education, salary, citizenship_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
`

func (e employeeRepo) CreateEmployee(ctx context.Context, worker entity.Employee) error {
	_, err := e.client.conn.Exec(
		ctx,
		createEmployeeQuery,
		worker.TinNum,
		worker.Snils,
		worker.PassportData,
		worker.PhoneNum,
		worker.FirstName,
		worker.LastName,
		worker.Patronymic,
		worker.Address,
		worker.Title.ID,
		worker.HiringDate,
		worker.Unit.ID,
		worker.Education,
		worker.Salary,
		worker.Citizenship.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return nil
}

const getTitlesQuery = `
SELECT id, name
FROM usecase.titles;
`

func (e employeeRepo) GetTitles(ctx context.Context) ([]entity.Title, error) {
	res, err := e.client.conn.Query(ctx, getTitlesQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	titles := make([]entity.Title, 0, 100)
	for res.Next() {
		var title entity.Title
		if err := res.Scan(&title.ID, &title.Name); err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		titles = append(titles, title)
	}
	return titles, nil
}

const getCitizenshipQuery = `
SELECT id, name
FROM usecase.citizenships;
`

func (e employeeRepo) GetCitizenships(ctx context.Context) ([]entity.Citizenship, error) {
	res, err := e.client.conn.Query(ctx, getTitlesQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	citizenships := make([]entity.Citizenship, 0, 100)
	for res.Next() {
		var citizenship entity.Citizenship
		if err := res.Scan(&citizenship.ID, &citizenship.Name); err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		citizenships = append(citizenships, citizenship)
	}
	return citizenships, nil
}

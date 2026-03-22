package postgres

import (
	"context"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
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

const isEmployeeExistsQuery = `
SELECT COUNT(*)
FROM usecase.employees
WHERE tin_num=$1 OR snils_num=$2 OR passport_data=$3 OR phone_num=$4;
`

func (e employeeRepo) IsEmployeeExists(ctx context.Context, worker entity.Employee) error {
	res, err := e.client.conn.Query(
		ctx,
		isEmployeeExistsQuery,
		worker.TinNum,
		worker.Snils,
		worker.PassportData,
		worker.PhoneNum,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if res.Next() {
		return persistent.ErrQueryExec
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}

	if count == 0 {
		return persistent.ErrEmployeeNotFound
	}
	return nil
}

const createEmployeeQuery = `
INSERT INTO usecase.employees (
	tin_num, 
	snils_num, 
	passport_data, 
	phone_num, 
	first_name, 
	last_name, 
	patronymic, 
	address, 
	title_id, 
	hiring_date, 
	unit_id, 
	education, 
	salary, 
	citizenship_id
)
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

const updateEmployeeQuery = `
UPDATE usecase.employee
SET tin_num=$1,
	snils_num=$2, 
	passport_data=$3, 
	phone_num=$4, 
	first_name=$5, 
	last_name=$6, 
	patronymic=$7, 
	address=$8, 
	title_id=$9, 
	hiring_date=$10, 
	unit_id=$11, 
	education=$12, 
	salary=$13, 
	citizenship_id=$14
WHERE id=$15;
`

func (e employeeRepo) UpdateEmployee(ctx context.Context, worker entity.Employee) error {
	_, err := e.client.conn.Exec(
		ctx,
		updateEmployeeQuery,
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
		worker.EmployeeID,
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

const countEmployeeWithCitizenships = `
SELECT usecase.citizenships.id, usecase.citizenships.name, COUNT(*)
FROM usecase.employees
	JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
GROUP BY usecase.citizenships.id, usecase.citizenships.name;
`

func (e employeeRepo) CountEmployeeWithCitizenships(
	ctx context.Context,
) ([]employee.EmployeeCitizenshipStatistic, error) {
	res, err := e.client.conn.Query(ctx, countEmployeeWithCitizenships)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	stats := make([]employee.EmployeeCitizenshipStatistic, 0, 100)
	for res.Next() {
		var employeeStat employee.EmployeeCitizenshipStatistic

		err := res.Scan(
			&employeeStat.Citizenship.ID,
			&employeeStat.Citizenship.Name,
			&employeeStat.EmployeeCount,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
		stats = append(stats, employeeStat)
	}
	return stats, nil
}

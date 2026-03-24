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

const countEmployeesWithCitizenship = `
SELECT usecase.citizenships.id, usecase.citizenships.name, COUNT(*)
FROM usecase.employees
	JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
GROUP BY usecase.citizenships.id, usecase.citizenships.name;
`

func (e employeeRepo) CountEmployeesWithCitizenship(
	ctx context.Context,
) ([]employee.EmployeeCitizenshipStatistic, error) {
	res, err := e.client.conn.Query(ctx, countEmployeesWithCitizenship)
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

const countEmployeesWithSalaryBoundsQuery = `
SELECT COUNT(*)
FROM usecase.employees
WHERE salary >= $1 AND salary <= $2 AND title_id=$3;
`

func (e employeeRepo) CountEmployeesWithSalaryBounds(
	ctx context.Context,
	titleID int,
	bounds employee.SalaryBounds,
) (int, error) {
	res, err := e.client.conn.Query(
		ctx,
		countEmployeesWithSalaryBoundsQuery,
		bounds.DownBoundary,
		bounds.UpBoundary,
		titleID,
	)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	if !res.Next() {
		return 0, persistent.ErrQueryExec
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	return count, nil
}

const getEmployeesByNameQuery = `
SELECT
	usecase.employees.id,
	usecase.employees.tin_num, 
	usecase.employees.snils_num, 
	usecase.employees.passport_data, 
	usecase.employees.phone_num, 
	usecase.employees.first_name, 
	usecase.employees.last_name, 
	usecase.employees.patronymic, 
	usecase.employees.address, 
	usecase.employees.title_id, 
	usecase.titles.name,
	usecase.employees.hiring_date, 
	usecase.employees.unit_id,
	usecase.divisions.name,
	usecase.divisions.type,
	usecase.divisions.state_size,
	usecase.divisions.superdivision_id,
	usecase.employees.education, 
	usecase.employees.salary, 
	usecase.employees.citizenship_id,
	usecase.citizenships.name
FROM 
	usecase.employees
		JOIN
	usecase.divisions
	ON usecase.employees.unit_id=usecase.divsions.unit_id
		JOIN
	usecase.titles
	ON usecase.employees.title_id=usecase.titles.id
		JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
WHERE 
	usecase.employees.first_name LIKE '%$1%' AND
	usecase.employees.last_name LIKE '%$2%' AND
	usecase.employees.patronymic LIKE '%$3%'
OFFSET $4
LIMIT $5;
`

func (e employeeRepo) GetEmployeesByName(ctx context.Context, filter employee.NamesFilter) ([]entity.Employee, error) {
	res, err := e.client.conn.Query(
		ctx,
		getEmployeesByNameQuery,
		filter.FirstName,
		filter.LastName,
		filter.Patronymic,
		filter.PageNum*employee.PaginationSize,
		employee.PaginationSize,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]entity.Employee, 0, 1_000)
	for res.Next() {
		var worker entity.Employee
		err := res.Scan(
			&worker.EmployeeID,
			&worker.TinNum,
			&worker.Snils,
			&worker.PassportData,
			&worker.PhoneNum,
			&worker.FirstName,
			&worker.LastName,
			&worker.Patronymic,
			&worker.Address,
			&worker.Title.ID,
			&worker.Title.Name,
			&worker.HiringDate,
			&worker.Unit.ID,
			&worker.Unit.Name,
			&worker.Unit.Type,
			&worker.Unit.StateSize,
			&worker.Unit.SuperdivisionID,
			&worker.Education,
			&worker.Salary,
			&worker.Citizenship.ID,
			&worker.Citizenship.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
	}
	return list, nil
}

const getEmployeesByPassportDataQuery = `
SELECT
	usecase.employees.id,
	usecase.employees.tin_num, 
	usecase.employees.snils_num, 
	usecase.employees.passport_data, 
	usecase.employees.phone_num, 
	usecase.employees.first_name, 
	usecase.employees.last_name, 
	usecase.employees.patronymic, 
	usecase.employees.address, 
	usecase.employees.title_id, 
	usecase.titles.name,
	usecase.employees.hiring_date, 
	usecase.employees.unit_id,
	usecase.divisions.name,
	usecase.divisions.type,
	usecase.divisions.state_size,
	usecase.divisions.superdivision_id,
	usecase.employees.education, 
	usecase.employees.salary, 
	usecase.employees.citizenship_id,
	usecase.citizenships.name
FROM 
	usecase.employees
		JOIN
	usecase.divisions
	ON usecase.employees.unit_id=usecase.divsions.unit_id
		JOIN
	usecase.titles
	ON usecase.employees.title_id=usecase.titles.id
		JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
WHERE 
	usecase.employees.passport_data LIKE '%$1%'
OFFSET $2
LIMIT $3;
`

func (e employeeRepo) GetEmployeesByPassportData(
	ctx context.Context,
	filter employee.PassportFilter,
) ([]entity.Employee, error) {
	res, err := e.client.conn.Query(
		ctx,
		getEmployeesByNameQuery,
		filter.PassportData,
		filter.PageNum*employee.PaginationSize,
		employee.PaginationSize,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]entity.Employee, 0, 1_000)
	for res.Next() {
		var worker entity.Employee
		err := res.Scan(
			&worker.EmployeeID,
			&worker.TinNum,
			&worker.Snils,
			&worker.PassportData,
			&worker.PhoneNum,
			&worker.FirstName,
			&worker.LastName,
			&worker.Patronymic,
			&worker.Address,
			&worker.Title.ID,
			&worker.Title.Name,
			&worker.HiringDate,
			&worker.Unit.ID,
			&worker.Unit.Name,
			&worker.Unit.Type,
			&worker.Unit.StateSize,
			&worker.Unit.SuperdivisionID,
			&worker.Education,
			&worker.Salary,
			&worker.Citizenship.ID,
			&worker.Citizenship.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
	}
	return list, nil
}

const getEmployeesByNameInDivisionQuery = `
SELECT
	usecase.employees.id,
	usecase.employees.tin_num, 
	usecase.employees.snils_num, 
	usecase.employees.passport_data, 
	usecase.employees.phone_num, 
	usecase.employees.first_name, 
	usecase.employees.last_name, 
	usecase.employees.patronymic, 
	usecase.employees.address, 
	usecase.employees.title_id, 
	usecase.titles.name,
	usecase.employees.hiring_date, 
	usecase.employees.unit_id,
	usecase.divisions.name,
	usecase.divisions.type,
	usecase.divisions.state_size,
	usecase.divisions.superdivision_id,
	usecase.employees.education, 
	usecase.employees.salary, 
	usecase.employees.citizenship_id,
	usecase.citizenships.name
FROM 
	usecase.employees
		JOIN
	usecase.divisions
	ON usecase.employees.unit_id=usecase.divsions.unit_id
		JOIN
	usecase.titles
	ON usecase.employees.title_id=usecase.titles.id
		JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
WHERE 
	usecase.employees.first_name LIKE '%$1%' AND
	usecase.employees.last_name LIKE '%$2%' AND
	usecase.employees.patronymic LIKE '%$3%' AND
	usecase.divisions.name LIKE '%$4%' AND
	usecase.divisions.type LIKE '%$5%'
OFFSET $6
LIMIT $7;
`

func (e employeeRepo) GetEmployeesByNameInDivision(
	ctx context.Context,
	filter employee.NamesFilter,
	div employee.UnitFilter,
) ([]entity.Employee, error) {
	res, err := e.client.conn.Query(
		ctx,
		getEmployeesByNameQuery,
		filter.FirstName,
		filter.LastName,
		filter.Patronymic,
		div.Name,
		div.Type,
		filter.PageNum*employee.PaginationSize,
		employee.PaginationSize,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]entity.Employee, 0, 1_000)
	for res.Next() {
		var worker entity.Employee
		err := res.Scan(
			&worker.EmployeeID,
			&worker.TinNum,
			&worker.Snils,
			&worker.PassportData,
			&worker.PhoneNum,
			&worker.FirstName,
			&worker.LastName,
			&worker.Patronymic,
			&worker.Address,
			&worker.Title.ID,
			&worker.Title.Name,
			&worker.HiringDate,
			&worker.Unit.ID,
			&worker.Unit.Name,
			&worker.Unit.Type,
			&worker.Unit.StateSize,
			&worker.Unit.SuperdivisionID,
			&worker.Education,
			&worker.Salary,
			&worker.Citizenship.ID,
			&worker.Citizenship.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
	}
	return list, nil
}

const getEmployeesByPassportDataInDivisionQuery = `
SELECT
	usecase.employees.id,
	usecase.employees.tin_num, 
	usecase.employees.snils_num, 
	usecase.employees.passport_data, 
	usecase.employees.phone_num, 
	usecase.employees.first_name, 
	usecase.employees.last_name, 
	usecase.employees.patronymic, 
	usecase.employees.address, 
	usecase.employees.title_id, 
	usecase.titles.name,
	usecase.employees.hiring_date, 
	usecase.employees.unit_id,
	usecase.divisions.name,
	usecase.divisions.type,
	usecase.divisions.state_size,
	usecase.divisions.superdivision_id,
	usecase.employees.education, 
	usecase.employees.salary, 
	usecase.employees.citizenship_id,
	usecase.citizenships.name
FROM 
	usecase.employees
		JOIN
	usecase.divisions
	ON usecase.employees.unit_id=usecase.divsions.unit_id
		JOIN
	usecase.titles
	ON usecase.employees.title_id=usecase.titles.id
		JOIN
	usecase.citizenships
	ON usecase.employees.citizenship_id=usecase.citizenships.id
WHERE 
	usecase.employees.passport_data LIKE '%$1%' AND
	usecase.divisions.name LIKE '%$2%' AND
	usecase.divisions.type LIKE '%$3%'
OFFSET $4
LIMIT $5;
`

func (e employeeRepo) GetEmployeesByPassportDataInDivision(
	ctx context.Context,
	filter employee.PassportFilter,
	div employee.UnitFilter,
) ([]entity.Employee, error) {
	res, err := e.client.conn.Query(
		ctx,
		getEmployeesByNameQuery,
		filter.PassportData,
		div.Name,
		div.Type,
		filter.PageNum*employee.PaginationSize,
		employee.PaginationSize,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
	}
	defer res.Close()

	list := make([]entity.Employee, 0, 1_000)
	for res.Next() {
		var worker entity.Employee
		err := res.Scan(
			&worker.EmployeeID,
			&worker.TinNum,
			&worker.Snils,
			&worker.PassportData,
			&worker.PhoneNum,
			&worker.FirstName,
			&worker.LastName,
			&worker.Patronymic,
			&worker.Address,
			&worker.Title.ID,
			&worker.Title.Name,
			&worker.HiringDate,
			&worker.Unit.ID,
			&worker.Unit.Name,
			&worker.Unit.Type,
			&worker.Unit.StateSize,
			&worker.Unit.SuperdivisionID,
			&worker.Education,
			&worker.Salary,
			&worker.Citizenship.ID,
			&worker.Citizenship.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", persistent.ErrQueryExec, err)
		}
	}
	return list, nil
}

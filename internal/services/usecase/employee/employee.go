package employee

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
)

type Interactor struct {
	workerRepo IEmployeeRepo
	reporter   reportManager
}

func NewInteractor(workerRepo IEmployeeRepo) Interactor {
	return Interactor{
		workerRepo: workerRepo,
	}
}

func (e Interactor) CreateEmployee(ctx context.Context, employee entity.Employee) error {
	err := e.workerRepo.IsEmployeeExists(ctx, employee)
	if err == nil {
		return ErrEmployeeExists
	} else if !errors.Is(err, persistent.ErrEmployeeNotFound) {
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	err = e.CreateEmployee(ctx, employee)
	if err != nil {
		retErr := fmt.Errorf("%w: %s", ErrRepoInteract, err)
		if errors.Is(err, persistent.ErrCitizenshipNotFound) {
			retErr = fmt.Errorf("%w: %s", ErrCitizenshipNotFound, err)
		} else if errors.Is(err, persistent.ErrTitleNotFound) {
			retErr = fmt.Errorf("%w: %s", ErrTitleNotFound, err)
		}
		return retErr
	}
	return nil
}

func (e Interactor) UpdateEmployee(ctx context.Context, employee entity.Employee) error {
	err := e.workerRepo.UpdateEmployee(ctx, employee)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return nil
}

func (e Interactor) GetTitles(ctx context.Context) ([]entity.Title, error) {
	titles, err := e.workerRepo.GetTitles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return titles, nil
}

func (e Interactor) GetCitizenships(ctx context.Context) ([]entity.Citizenship, error) {
	citizenships, err := e.GetCitizenships(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return citizenships, nil
}

func (e Interactor) CountEmployeesWithCitizenship(ctx context.Context) ([]EmployeeCitizenshipStatistic, error) {
	stats, err := e.workerRepo.CountEmployeesWithCitizenship(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return stats, nil
}

func (e Interactor) CountEmployeesWithSalaryBounds(
	ctx context.Context,
	titleID int,
	bounds SalaryBounds,
) (int, error) {
	count, err := e.workerRepo.CountEmployeesWithSalaryBounds(ctx, titleID, bounds)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return count, nil
}

func (e Interactor) GetEmployeesWithFilters(ctx context.Context, filters Filter, pageNum int) ([]entity.Employee, error) {
	var (
		list []entity.Employee
		err  error
	)

	if filters.Names.IsActive {
		if filters.Unit.IsActive {
			list, err = e.workerRepo.GetEmployeesByNameInDivision(ctx, filters.Names, filters.Unit)
		} else {
			list, err = e.workerRepo.GetEmployeesByName(ctx, filters.Names)
		}
	} else if filters.Passport.IsActive {
		if filters.Unit.IsActive {
			list, err = e.workerRepo.GetEmployeesByPassportDataInDivision(ctx, filters.Passport, filters.Unit)
		} else {
			list, err = e.workerRepo.GetEmployeesByPassportData(ctx, filters.Passport)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return list, nil
}

func (e Interactor) DeleteEmployee(ctx context.Context, employeeID int) (string, error) {
	worker, err := e.workerRepo.DeleteEmployee(ctx, employeeID)
	if err != nil {
		retErr := ErrRepoInteract
		if errors.Is(err, persistent.ErrEmployeeNotFound) {
			retErr = ErrEmployeeNotFound
		}
		return "", fmt.Errorf("%w: %s", retErr, err)
	}

	reportID, err := e.reporter.createDeletedEmployeeReport(worker)
	if err != nil {
		return "", err
	}
	return reportID, nil
}

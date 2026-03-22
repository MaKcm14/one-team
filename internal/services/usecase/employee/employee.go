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

package root

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Interactor struct {
	repo IRootRepo
}

func NewInteractor(repo IRootRepo) Interactor {
	return Interactor{
		repo: repo,
	}
}

func (r Interactor) GetUsers(ctx context.Context, filters user.Filters) ([]UserDTO, error) {
	list, err := r.repo.GetUsersByLogin(ctx, filters.LoginFilter)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	res := make([]UserDTO, 0, len(list))
	for _, val := range list {
		res = append(res, UserDTO{
			Login: val.User.Login,
			Role:  val.Role,
		})
	}
	return res, nil
}

func (r Interactor) GetRoles(ctx context.Context) ([]Role, error) {
	list, err := r.repo.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return list, nil
}

func (r Interactor) DeleteUser(ctx context.Context, login string) error {
	_, err := r.repo.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, persistent.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	role, err := r.repo.GetUserRole(ctx, login)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	if role == entity.AdminRole {
		return ErrUnableToDeleteAdmin
	}

	err = r.repo.DeleteUser(ctx, login)
	if err != nil {
		if errors.Is(err, persistent.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return nil
}

func (r Interactor) UpdateUserRole(ctx context.Context, dto UserDTO) error {
	role, err := r.repo.GetUserRole(ctx, dto.Login)
	if err != nil {
		if errors.Is(err, persistent.ErrRoleNotAssign) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	if role == entity.AdminRole {
		return ErrUnableToChangeAdminRole
	}

	err = r.repo.IsRoleExists(ctx, dto.Role)
	if err != nil {
		if errors.Is(err, persistent.ErrRoleNotFound) {
			return ErrRoleNotFound
		}
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}

	err = r.repo.UpdateUserRole(ctx, dto)
	if err != nil {
		if errors.Is(err, persistent.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%w: %s", ErrRepoInteract, err)
	}
	return nil
}

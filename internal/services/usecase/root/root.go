package root

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo IRootRepo
}

func NewInteractor(repo IRootRepo) Interactor {
	return Interactor{
		repo: repo,
	}
}

func (r Interactor) GetUsers(ctx context.Context) ([]UserDTO, error) {
	list, err := r.repo.GetUsers(ctx)
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

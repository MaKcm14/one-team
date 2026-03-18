package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Interactor struct {
	log      *slog.Logger
	cfg      config.AuthConfig
	authRepo user.IAuthRepo
}

func NewInteractor(
	log *slog.Logger,
	cfg config.AuthConfig,
	authRepo user.IAuthRepo,
) Interactor {
	return Interactor{
		log:      log,
		cfg:      cfg,
		authRepo: authRepo,
	}
}

func (auth Interactor) Login(ctx context.Context, creds user.Credentials) (user.UserDTO, error) {
	userInfo, err := auth.authRepo.GetUser(ctx, creds.Login)
	if err != nil {
		retErr := fmt.Errorf("%w: %s", user.ErrRepoInteract, err)
		if errors.Is(err, persistent.ErrUserNotFound) {
			retErr = user.ErrUserNotFound
		}
		return user.UserDTO{}, retErr
	}

	err = auth.checkPassword(userInfo.HashPWD, creds.Password+fmt.Sprint(userInfo.Salt+auth.cfg.GlobalPwdSalt))
	if err != nil {
		return user.UserDTO{}, user.ErrWrongPassword
	}

	role, err := auth.authRepo.GetUserRole(ctx, creds.Login)
	if err != nil {
		retErr := fmt.Errorf("%w: %s", user.ErrRepoInteract, err)
		if errors.Is(err, persistent.ErrRoleNotAssign) {
			retErr = user.ErrRoleNotAssign
		}
		return user.UserDTO{}, retErr
	}
	return user.UserDTO{
		User: userInfo,
		Role: role,
	}, nil
}

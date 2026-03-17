package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/one-team/internal/repository/persistent"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Interactor struct {
	log      *slog.Logger
	authRepo user.IAuthRepo
}

func NewInteractor(log *slog.Logger, authRepo user.IAuthRepo) Interactor {
	return Interactor{
		log:      log,
		authRepo: authRepo,
	}
}

func (auth Interactor) Login(ctx context.Context, creds user.Credentials) error {
	userInfo, err := auth.authRepo.GetUser(ctx, creds.Login)
	if err != nil {
		if errors.Is(err, persistent.ErrUserNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("%w: %s", user.ErrRepoInteract, err)
	}

	err = auth.checkPassword(userInfo.HashPWD, creds.Password)
	if err != nil {
		return user.ErrWrongPassword
	}
	return nil
}

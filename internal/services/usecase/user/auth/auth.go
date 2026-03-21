package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/one-team/internal/config"
	entity "github.com/MaKcm14/one-team/internal/entity/user"
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

func (auth Interactor) SignUp(ctx context.Context, dto user.UserSignUpDTO) error {
	_, err := auth.authRepo.GetUser(ctx, dto.Creds.Login)
	if errors.Is(err, persistent.ErrUserNotFound) {
		if err := auth.verifyPassword(dto.Creds.Password); err != nil {
			return fmt.Errorf("%w: %w", user.ErrSignUp, err)
		}

		userSalt := auth.generateSalt()
		hashPwd, err := auth.hashPassword(dto.Creds.Password, userSalt)
		if err != nil {
			return fmt.Errorf("%w: %s", user.ErrHashPassword, err)
		}

		err = auth.authRepo.CreateUser(ctx, user.UserDTO{
			User: entity.User{
				Login:   dto.Creds.Login,
				HashPWD: string(hashPwd),
			},
			Role: dto.Role,
		})
		if err != nil {
			if errors.Is(err, persistent.ErrRoleNotFound) {
				return fmt.Errorf("%w: %w: %s", user.ErrSignUp, user.ErrRoleNotFound, err)
			}
			return fmt.Errorf("%w: %s", user.ErrRepoInteract, err)
		}
		return nil

	} else if err != nil {
		return fmt.Errorf("%w: %s", user.ErrRepoInteract, err)
	}
	return fmt.Errorf("%w: %w", user.ErrSignUp, user.ErrUserAlreadyExists)
}

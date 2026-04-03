package root

import (
	"context"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type IRootRepo interface {
	GetUsers(ctx context.Context) ([]user.UserDTO, error)
	GetRoles(ctx context.Context) ([]Role, error)
}

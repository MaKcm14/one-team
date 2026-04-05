package root

import (
	"context"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
)

type IRootRepoWriter interface {
	DeleteUser(ctx context.Context, login string) error
	UpdateUserRole(ctx context.Context, user UserDTO) error
}

type IRootRepoReader interface {
	IsRoleExists(ctx context.Context, role entity.Role) error
	GetUser(ctx context.Context, login string) (entity.User, error)
	GetUsersByLogin(ctx context.Context, filter user.LoginFilter) ([]user.UserDTO, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetUserRole(ctx context.Context, login string) (entity.Role, error)
}

type IRootRepo interface {
	IRootRepoReader
	IRootRepoWriter
}

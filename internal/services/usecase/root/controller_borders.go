package root

import (
	"context"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type IRootServiceReader interface {
	GetUsers(ctx context.Context, filters user.Filters) ([]UserDTO, error)
	GetRoles(ctx context.Context) ([]Role, error)
}

type IRootServiceWriter interface {
	DeleteUser(ctx context.Context, login string) error
	UpdateUserRole(ctx context.Context, user UserDTO) error
}

type IRootService interface {
	IRootServiceReader
	IRootServiceWriter
}

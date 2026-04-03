package root

import "context"

type IRootServiceReader interface {
	GetUsers(ctx context.Context) ([]UserDTO, error)
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

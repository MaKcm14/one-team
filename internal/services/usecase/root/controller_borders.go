package root

import "context"

type IRootService interface {
	GetUsers(ctx context.Context) ([]UserDTO, error)
	GetRoles(ctx context.Context) ([]Role, error)
	DeleteUser(ctx context.Context, login string) error
}

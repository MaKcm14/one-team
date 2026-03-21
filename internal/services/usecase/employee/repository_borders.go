package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type IEmployeeRepoWriter interface {
	CreateEmployee(ctx context.Context, worker entity.Employee) error
}

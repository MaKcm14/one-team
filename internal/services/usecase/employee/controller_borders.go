package employee

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type IEmployeeServiceModifier interface {
	CreateEmployee(ctx context.Context, employee entity.Employee) error
}

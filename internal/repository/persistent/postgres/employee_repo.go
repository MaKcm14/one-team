package postgres

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/employee"
)

type employeeRepo struct {
	client *postgresClient
}

func (e employeeRepo) CreateEmployee(ctx context.Context, worker entity.Employee) error {
	return nil
}

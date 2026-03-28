package postgres

import (
	"fmt"

	"github.com/MaKcm14/one-team/internal/app/logger"
	"github.com/MaKcm14/one-team/internal/config"
)

type Repository struct {
	userRepo
	employeeRepo
	divisionRepo

	client *postgresClient
}

func NewRepository(log logger.Logger, cfg config.DBConfig) (Repository, error) {
	client, err := newPostgresClient(log, cfg)
	if err != nil {
		return Repository{}, err
	}
	return Repository{
		client: client,
		userRepo: userRepo{
			client: client,
		},
		employeeRepo: employeeRepo{
			client: client,
		},
		divisionRepo: divisionRepo{
			client: client,
		},
	}, nil
}

func as(str string) string {
	return fmt.Sprintf("%%%s%%", str)
}

func (r Repository) Close() {
	r.client.close()
}

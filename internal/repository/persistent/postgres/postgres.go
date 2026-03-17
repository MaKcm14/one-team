package postgres

import (
	"log/slog"

	"github.com/MaKcm14/one-team/internal/config"
)

type Repository struct {
	userRepo

	client *postgresClient
}

func NewRepository(log *slog.Logger, cfg config.DBConfig) (Repository, error) {
	client, err := newPostgresClient(log, cfg)
	if err != nil {
		return Repository{}, err
	}
	return Repository{
		client: client,
		userRepo: userRepo{
			client: client,
		},
	}, nil
}

func (r Repository) Close() {
	r.client.close()
}

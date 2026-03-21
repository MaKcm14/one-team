package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MaKcm14/one-team/internal/app/logger"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/repository/persistent"
)

type postgresClient struct {
	log  logger.Logger
	conn *pgxpool.Pool
}

func newPostgresClient(log logger.Logger, cfg config.DBConfig) (*postgresClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrConnWithDB, err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %s", persistent.ErrConnWithDB, err)
	}
	return &postgresClient{
		log:  log,
		conn: conn,
	}, nil
}

func (p *postgresClient) close() {
	p.conn.Close()
}

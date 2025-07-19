package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/subscriptions_api/internal/logger"
)

var PostgresDB *pgx.Conn

func NewPostgresDB(ctx context.Context, dsn string) (*pgx.Conn, error) {

	logger.L.Debug("Start connecting to Postgres DB")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("[NewPostgresDB|connect DB] , %w", err)
	}

	// Проверка подключения
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("[NewPostgresDB|ping] ,%w", err)
	}

	PostgresDB = conn
	logger.L.Info("Successful connected to DB")
	return PostgresDB, nil
}

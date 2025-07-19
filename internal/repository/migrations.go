package repository

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	p "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/subscriptions_api/internal/logger"
)

func RunMigrations(ctx context.Context, conn *p.Conn) error {

	logger.L.Debug("Start processing migrations")
	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	// Создаем драйвер для migrate напрямую из pgx.Conn
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return err
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к миграциям
		"pgx", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.L.Debug("Applied migrations")
	return nil
}

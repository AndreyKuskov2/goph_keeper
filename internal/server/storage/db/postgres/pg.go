package postgres

import (
	"context"
	"fmt"
	"goph_keeper/internal/server/config"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg *config.Database) (*Postgres, error) {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err := initialMigrate(pool, cfg.MigrationsPath, cfg.Name); err != nil {
		return nil, err
	}

	return &Postgres{
		DB: pool,
	}, nil
}

func initialMigrate(pool *pgxpool.Pool, migrationPath, dbName string) error {
	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		dbName, driver)
	if err != nil {
		return fmt.Errorf("cannot create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cannot to apply migrations: %v", err)
	} else if err == migrate.ErrNoChange {
		slog.Debug("no new migrations. current database version is used")
	}
	slog.Debug("apply all new migrations")
	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	slog.Debug("ping pg")
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	slog.Debug("close pg connection")
	// FIX: После вызова Close поток выполнения блокируется. Понять как нормально закрывать соединение
	// pg.DB.Close()
}

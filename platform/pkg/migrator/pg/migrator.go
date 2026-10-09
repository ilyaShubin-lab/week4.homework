package pg

import (
	"boilerplates/platform/pkg/migrator"
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

var _ migrator.Migrator = (*Migrator)(nil)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, m.db, m.migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

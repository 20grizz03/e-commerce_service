package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrator структура для применения миграций.
type Migrator struct {
	srcDriver source.Driver // Драйвер источника миграций.
}

// MustGetNewMigrator создает новый экземпляр Migrator с встроенными SQL-файлами миграций.
// В случае ошибки вызывает panic.
func MustGetNewMigrator(sqlFiles embed.FS, dirName string) *Migrator {
	// Создаем новый драйвер источника миграций с встроенными SQL-файлами.
	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		panic(err)
	}
	return &Migrator{
		srcDriver: d,
	}
}

// ApplyMigrations применяет миграции к базе данных.
func (m *Migrator) ApplyMigrations(db *sql.DB) error {
	// Создаем экземпляр драйвера базы данных для PostgreSQL.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	// Создаем новый экземпляр мигратора с использованием драйвера источника и драйвера базы данных PostgreSQL.
	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	// Закрываем мигратор в конце работы функции.
	defer func() {
		migrator.Close()
	}()

	// Применяем миграции.
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations %v", err)
	}

	return nil
}

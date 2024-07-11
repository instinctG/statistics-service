package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
)

// MigrateDB выполняет миграцию базы данных, используя миграции из папки migrations.
func (d *Database) MigrateDB() error {
	fmt.Println("migrating our database")

	// Преобразование подключения pgxpool к *sql.DB
	db := stdlib.OpenDBFromPool(d.Client)
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	fmt.Println("successfully migrated the database")
	return nil
}

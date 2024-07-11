package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type Database struct {
	Client *pgxpool.Pool
}

// NewDatabase создает новое подключение к базе данных на основе конфигурации из переменных окружения.
func NewDatabase() (*Database, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{pool}, nil

}

// Ping проверяет соединение с базой данных.
func (d *Database) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx)
}

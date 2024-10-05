package database

import (
	"context"
	"fmt"

	"github.com/amedoeyes/hadath/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect() error {
	cfg := config.Get()
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser(), cfg.DBPassword(), cfg.DBHost(), cfg.DBPort(), cfg.DBName())
	var err error
	pool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func Disconnect() {
	pool.Close()
}

func Get() *pgxpool.Pool {
	return pool
}

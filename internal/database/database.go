package database

import (
	"context"
	"fmt"

	"github.com/amedoeyes/hadath/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func Connect(ctx context.Context) error {
	cfg := config.Get()
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser(), cfg.DBPassword(), cfg.DBHost(), cfg.DBPort(), cfg.DBName())
	var err error
	db, err = pgxpool.New(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func Disconnect() {
	db.Close()
}

func Get() *pgxpool.Pool {
	return db
}

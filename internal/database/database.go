package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/amedoeyes/hadath/config"
	"github.com/amedoeyes/hadath/internal/utility"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var db *pgxpool.Pool

func Connect(ctx context.Context) error {
	cfg := config.Get()
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	db, err = pgxpool.New(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("Connected to database %s at %s:%d", cfg.DBName, cfg.DBHost, cfg.DBPort)

	return nil
}

func Disconnect() {
	if db != nil {
		db.Close()
		log.Println("Disconnected from database")
	}
}

func MigrateUp() error {
	cfg := config.Get()
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	sqlDB, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	root, err := utility.FindProjectRoot()
	if err != nil {
		return err
	}

	if err := goose.Up(sqlDB, root+"/internal/database/migration"); err != nil {
		return err
	}

	return nil
}

func MigrateDown() error {
	cfg := config.Get()
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	sqlDB, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	root, err := utility.FindProjectRoot()
	if err != nil {
		return err
	}

	if err := goose.Reset(sqlDB, root+"/internal/database/migration"); err != nil {
		return err
	}

	return nil
}

func Get() *pgxpool.Pool {
	return db
}

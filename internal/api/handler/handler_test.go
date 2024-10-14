package handler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/amedoeyes/hadath/config"
	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/validator"
)

func TestMain(m *testing.M) {
	err := config.Load(".env.test")
	if err != nil {
		log.Fatal(err)
	}

	err = database.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = database.MigrateUp()
	if err != nil {
		log.Fatal(err)
	}

	validator.Init()

	code := m.Run()

	database.Disconnect()
	os.Exit(code)
}

func setupTest(t *testing.T) {
	t.Helper()
	_, err := database.Get().Exec(context.Background(), "TRUNCATE TABLE users, api_keys, events, bookings RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal(err)
	}
}

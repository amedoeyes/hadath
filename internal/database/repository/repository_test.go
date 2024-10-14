package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/amedoeyes/hadath/config"
	"github.com/amedoeyes/hadath/internal/database"
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

	code := m.Run()

	database.Disconnect()
	os.Exit(code)
}

func setupTest(t *testing.T) {
	t.Helper()
	err := database.MigrateUp()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := database.MigrateDown()
		if err != nil {
			t.Fatal(err)
		}
	})
}

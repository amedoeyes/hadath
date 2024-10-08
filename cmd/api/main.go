package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/amedoeyes/hadath/config"
	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/router"
)

func main() {
	config.Load()
	err := database.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.Disconnect()

	cfg := config.Get()
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost(), cfg.ServerPort())
	router := router.Setup()

	log.Printf("Starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amedoeyes/hadath/config"
	"github.com/amedoeyes/hadath/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return router
}

func main() {
	config.Load()
	err := database.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.Disconnect()

	cfg := config.Get()
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost(), cfg.ServerPort())
	log.Printf("Starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, RegisterRoutes()))
}

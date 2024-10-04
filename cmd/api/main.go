package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amedoeyes/hadath/config"
	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/headers", headers)
	return router
}

func main() {
	config.LoadConfig()
	cfg := config.Get()
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost(), cfg.ServerPort())
	log.Printf("Starting server at %s", addr)
	http.ListenAndServe(addr, RegisterRoutes())
}

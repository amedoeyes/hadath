package main

import (
	"fmt"
	"log"
	"net/http"

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
	log.Println("Starting server on :8000")
	http.ListenAndServe(":8000", RegisterRoutes())
}

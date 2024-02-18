package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/contest/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got /api/v1/contest/\n")
	})

	mux.HandleFunc("GET /api/v1/contest/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprint(w, "handling task with id=%v\n", id)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

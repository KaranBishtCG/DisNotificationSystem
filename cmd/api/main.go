package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := &http.Server{
		Addr:         ":8088",
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server started on port 8088")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Not Connected %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/config"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/db"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/handlers"
)

func main() {

	cfg := config.MustLoad()

	_, err := db.Connect(cfg.Database_url)
	if err != nil {
		log.Fatalf("db.Connect connection failed %v", err)
	}

	fmt.Println("DB Connection Established")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server started on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Not Connected %v", err)
	}
}

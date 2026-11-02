package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/auth"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/config"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/db"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/handlers"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/middleware"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/notifications"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.Connect(cfg.Database_url)
	if err != nil {
		log.Fatalf("db.Connect connection failed %v", err)
	}

	fmt.Println("DB Connection Established")

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	logger := slog.New(logHandler)

	nh := notifications.NewNotificationHandler(db, logger)
	ah := auth.NewAuthHandler(db, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("POST /notifications", nh.Create)
	mux.HandleFunc("POST /signup", ah.Signup)

	handler := middleware.RequestId(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server started on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Not Connected %v", err)
	}
}

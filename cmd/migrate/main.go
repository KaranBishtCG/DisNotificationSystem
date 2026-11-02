package main

import (
	"log"
	"os"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down>")
	}

	m, err := migrate.New("file://migrations", cfg.Database_url)
	if err != nil {
		log.Fatalf("migration.new : %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command %s", os.Args[1])
	}
}

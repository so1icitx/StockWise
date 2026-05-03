package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/so1icitx/StockWise/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	migrationsDir := flag.String("dir", "migrations", "directory containing Goose migration files")
	flag.Parse()

	command := "up"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set goose dialect: %v", err)
	}

	if err := run(command, db, *migrationsDir); err != nil {
		log.Fatalf("migration %q failed: %v", command, err)
	}
}

func run(command string, db *sql.DB, migrationsDir string) error {
	switch command {
	case "up":
		return goose.Up(db, migrationsDir)
	case "down":
		return goose.Down(db, migrationsDir)
	case "status":
		return goose.Status(db, migrationsDir)
	case "reset":
		return goose.Reset(db, migrationsDir)
	case "version":
		return goose.Version(db, migrationsDir)
	default:
		fmt.Fprintf(os.Stderr, "usage: go run ./cmd/migrate [-dir migrations] [up|down|status|reset|version]\n")
		return fmt.Errorf("unknown command %q", command)
	}
}

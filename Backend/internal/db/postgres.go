package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.MaxConnLifetime = time.Hour

	DB, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Connected to Supabase PostgreSQL")
}

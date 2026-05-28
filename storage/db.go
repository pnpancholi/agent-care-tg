package storage

import (
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("[agent-care-tg]: DB connection failed", err)
		os.Exit(1)
	}
	slog.Info("[agent-care-tg]: DB connected")
	return db
}

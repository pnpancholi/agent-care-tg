package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("[agent-care-tg]: DB connection failed", err)
	}
	log.Println("[agent-care-tg]: DB connected")
	return db
}

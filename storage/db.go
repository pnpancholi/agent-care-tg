package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("[agent-care-tg]: DB connection failed", err)
	}
	log.Println("[agent-care-tg]: DB connected")
	return db
}

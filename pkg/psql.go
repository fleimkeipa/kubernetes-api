package pkg

import (
	"os"

	"github.com/go-pg/pg"
)

func NewPSQLClient() *pg.DB {
	var opts = pg.Options{
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     os.Getenv("DB_ADDR"),
	}
	var db = pg.Connect(&opts)

	return db
}

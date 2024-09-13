package pkg

import (
	"fmt"
	"os"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPSQLClient() *pg.DB {
	var opts = pg.Options{
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     os.Getenv("DB_ADDR"),
	}
	var db = pg.Connect(&opts)

	if err := createSchema(db); err != nil {
		panic(err.Error())
	}

	return db
}

func createSchema(db *pg.DB) error {
	var models = []interface{}{
		(*model.Event)(nil),
	}

	for _, model := range models {
		var opts = &orm.CreateTableOptions{
			IfNotExists: true,
		}

		if err := db.Model(model).CreateTable(opts); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

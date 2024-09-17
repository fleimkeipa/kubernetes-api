package pkg

import (
	"fmt"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"
)

func NewPSQLClient() *pg.DB {
	var opts = pg.Options{
		Database: viper.GetString("database.name"),
		User:     viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Addr:     viper.GetString("database.addr"),
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
		(*model.User)(nil),
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

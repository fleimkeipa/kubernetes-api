package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var test_db *pg.DB

func init() {
	var dbName = "test_db"
	os.Setenv("DB_ADDR", "localhost:5432")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", dbName)

	// createTestDB()

	if err := createTestSchema(test_db); err != nil {
		log.Fatal(err)
	}

	test_db = initTestDBClient()
}

func initTestDBClient() *pg.DB {
	var opts = pg.Options{
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     os.Getenv("DB_ADDR"),
	}
	var db = pg.Connect(&opts)

	return db
}

func createTestDB() {
	var conninfo = fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create database " + os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
}

// createSchema creates database schema for Event
func createTestSchema(db *pg.DB) error {
	var models = []interface{}{
		(*model.Event)(nil),
		(*model.User)(nil),
	}

	for _, model := range models {
		var opts = orm.CreateTableOptions{
			Temp:        true,
			IfNotExists: true,
		}
		time.Sleep(time.Millisecond * 500)
		err := db.
			Model(model).
			CreateTable(&opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func addTempData(data interface{}) error {
	_, err := test_db.Model(data).Insert()
	if err != nil {
		return err
	}

	return nil
}

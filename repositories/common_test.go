package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/joho/godotenv"
)

var test_db *pg.DB

func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	log.Println(".env file loaded successfully")
}

func init() {
	loadEnv()

	var addr = strings.Split(os.Getenv("DB_ADDR"), ":")
	if len(addr) < 2 {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
	} else {
		os.Setenv("DB_HOST", addr[0])
		os.Setenv("DB_PORT", addr[1])
	}

	createTestDB()

	test_db = initTestDBClient()

	if err := createTestSchema(test_db); err != nil {
		log.Fatal(err)
	}
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

	var query = fmt.Sprintf(`
		DROP DATABASE IF EXISTS %s;
	`, "test_db")
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = fmt.Sprintf(`
		CREATE DATABASE %s;
	`, "test_db")
	_, err = db.Exec(query)
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

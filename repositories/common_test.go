package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/fleimkeipa/kubernetes-api/config"
	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"
)

var test_db *pg.DB

func loadEnv() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	log.Println(".env file loaded successfully")
}

func init() {
	loadEnv()

	var addr = strings.Split(viper.GetString("database.addr"), ":")
	if len(addr) < 2 {
		viper.Set("database.host", "localhost")
		viper.Set("database.port", "5432")
	} else {
		viper.Set("database.host", addr[0])
		viper.Set("database.port", addr[1])
	}

	createTestDB()

	test_db = initTestDBClient()

	if err := createTestSchema(test_db); err != nil {
		log.Fatal(err)
	}
}

func initTestDBClient() *pg.DB {
	var opts = pg.Options{
		Database: viper.GetString("database.name"),
		User:     viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Addr:     viper.GetString("database.addr"),
	}
	var db = pg.Connect(&opts)

	return db
}

func createTestDB() {
	var conninfo = fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
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

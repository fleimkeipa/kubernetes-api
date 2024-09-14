package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/lib/pq"
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

	createTestDB()

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
	models := []interface{}{
		(*model.Event)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestEventRepository_Create(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx   context.Context
		event *model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.Background(),
				event: &model.Event{
					Kind:         "pod",
					EventKind:    "create",
					CreationTime: time.Now(),
					Owner: model.User{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			want:    &model.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &EventRepository{
				db: tt.fields.db,
			}
			got, err := rc.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

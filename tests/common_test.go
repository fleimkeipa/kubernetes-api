package tests

import (
	"fmt"

	"github.com/go-pg/pg"
)

var (
	test_db     *pg.DB
	terminateDB = func() {}
)

func addTempData(data interface{}) error {
	_, err := test_db.Model(data).Insert()
	if err != nil {
		return err
	}

	return nil
}

func clearTable(tableName string) error {
	query := fmt.Sprintf("TRUNCATE %s; DELETE FROM %s", tableName, tableName)
	_, err := test_db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

func GetColumns(db *sql.DB, table string) (columns []string) {
	result := &Result{}
	row, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + table + "'")
	if err != nil {
		log.Panicf("Error: Can't get columns")
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&result.Field)
		if err != nil {
			log.Panicf("Error: Can't scan row")
		}
		columns = append(columns, result.Field)
	}
	return columns[:len(columns)-3]
}

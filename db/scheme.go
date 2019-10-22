package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Result struct {
	Field   sql.NullString
	Type    sql.NullString
	Null    sql.NullString
	Key     sql.NullString
	Default sql.NullString
	Extra   sql.NullString
}

func GetColumns(db *sql.DB, table string) (columns []string) {
	result := &Result{}
	row, err := db.Query("SHOW COLUMNS FROM " + table)
	if err != nil {
		log.Panicf("Error: Can't get columns")
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		if err != nil {
			log.Panicf("Error: Can't scan row")
		}
		columns = append(columns, result.Field.String)
	}
	return columns
}

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

func GetColumns(db *sql.DB, table string) (columns []string, err error) {
	result := &Result{}
	row, err := db.Query("SHOW COLUMNS FROM " + table)
	if err != nil {
		log.Printf("Error: Can't get columns %v", err)
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		if err != nil {
			log.Printf("Error: Can't scan row %v", err)
		}
		columns = append(columns, result.Field.String)
	}
	return columns, err
}

func GetTables(db *sql.DB) (tables []string, err error) {
	var table string
	res, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Printf("Error: Can't get tables %v", err)
	}
	for res.Next() {
		res.Scan(&table)
		tables = append(tables, table)
	}
	return tables, err
}

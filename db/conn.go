package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Conn(conn string) *sql.DB {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatalf("Error: Can't connect DB %v", err)
	}

	return db
}

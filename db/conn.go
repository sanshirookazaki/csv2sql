package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Conn(conn string) *sql.DB {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Panicf("Error: Can't connect DB")
	}

	return db
}

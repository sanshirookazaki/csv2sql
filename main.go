package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sanshirookazaki/csv2sql/csv"
	"github.com/sanshirookazaki/csv2sql/db"
)

var (
	database = flag.String("d", "", "Name of Database")
	host     = flag.String("h", "127.0.0.1", "Host of Database")
	port     = flag.Int("P", 3306, "Database port number")
	user     = flag.String("u", "root", "Database username")
	password = flag.String("p", "", "Database password")
)

func main() {
	flag.Parse()
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", *user, *password, *host, *port, *database)
	db := db.Conn(conn)
	defer db.Close()

	if len(os.Args) < 2 {
		log.Panicf("Error: CSV path is required")
	}
	dir := os.Args[len(os.Args)-1]
	csvs := csv.FindCsv(dir)
	tables := createTableList(csvs)
	tx, err := db.Begin()
	if err != nil {
		log.Panicf("Error: Can't begin transaction")
	}
	for i, csv := range csvs {
		mysql.RegisterLocalFile(csv)
		_, err := db.Exec("LOAD DATA LOCAL INFILE '" + csv + "' INTO TABLE " + tables[i] + " FIELDS TERMINATED BY ',' IGNORE 1 LINES")
		if err != nil {
			tx.Rollback()
			log.Panicf("Error: Query faild")
		}
	}

	tx.Commit()
}

func createTableList(paths []string) (tableList []string) {
	for _, path := range paths {
		tableSlice := strings.Split(path, "/")
		table := strings.Join(tableSlice[:len(tableSlice)-1], "_")
		tableList = append(tableList, table)
	}

	return tableList
}

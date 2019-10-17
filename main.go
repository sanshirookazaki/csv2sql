package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	specific = flag.String("S", "", "Import specific tables")
	separate = flag.Bool("s", false, "Separate CSV into 2 types")
	ignore   = flag.Bool("i", true, "Ignore 1st line of CSV")
)

func main() {
	flag.Parse()
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", *user, *password, *host, *port, *database)
	DB := db.Conn(conn)
	defer DB.Close()

	if len(os.Args) < 2 {
		log.Panicf("Error: CSV path is required")
	}

	dir := os.Args[len(os.Args)-1]
	basePath, _ := filepath.Abs(dir)

	csvs := csv.FindCsv(basePath)
	tables := createTableList(csvs, basePath)
	tx, err := DB.Begin()
	if err != nil {
		log.Panicf("Error: Can't begin transaction")
	}
	csvs = filterSpecificTables(csvs, *specific)
	tables = filterSpecificTables(tables, *specific)

	for i, csvPath := range csvs {
		mysql.RegisterLocalFile(csvPath)
		var err error

		dbColumns := db.GetColumns(DB, tables[i])
		csvColumns := csv.GetColumns(csvPath)
		fmt.Println(csvColumns, dbColumns)
		diffColumns := diffSlice(dbColumns, csvColumns)
		fmt.Println(diffColumns)

		query := "LOAD DATA LOCAL INFILE '" + csvPath + "' INTO TABLE " + tables[i] + " FIELDS TERMINATED BY ',' "
		if len(diffColumns) == 0 {
			if *ignore {
				_, err = DB.Exec(query + " IGNORE 1 LINES")
			} else {
				_, err = DB.Exec(query)
			}
		} else {
			csvFile := getFileNameWithoutExt(csvPath)
			sets := " SET "
			for i, column := range diffColumns {
				sets += column + " = " + csvFile + " "
				if i != (len(diffColumns) - 1) {
					sets += ", "
				}
			}

			fmt.Println(sets)
			if *ignore {
				_, err = DB.Exec(query + " IGNORE 1 LINES " + sets)
			} else {
				_, err = DB.Exec(query + sets)
			}
		}

		if err != nil {
			tx.Rollback()
			fmt.Println(csvPath, "->", tables[i])
			log.Panicf("Error: Query faild")
		}
		fmt.Println(csvPath, "import to", tables[i])
	}

	tx.Commit()
	fmt.Println("Complete !!")
}

func createTableList(paths []string, basePath string) (tableList []string) {
	for _, path := range paths {
		rpath, _ := filepath.Rel(basePath, path)

		pathSlice := strings.Split(rpath, "/")
		var tableParts []string
		// the first charactor in filename is number
		if *separate && !initialIsInt(pathSlice[len(pathSlice)-1]) {
			pathSlice[len(pathSlice)-1] = strings.TrimRight(pathSlice[len(pathSlice)-1], ".csv")
			tableParts = pathSlice
		} else {
			tableParts = pathSlice[:len(pathSlice)-1]
		}
		table := strings.Join(tableParts, "_")
		tableList = append(tableList, table)
	}

	return tableList
}

func filterSpecificTables(list []string, specific string) (result []string) {
	for _, e := range list {
		if strings.Contains(e, specific) {
			result = append(result, e)
		}
	}
	return result
}

func initialIsInt(s string) bool {
	_, err := strconv.Atoi(getInitial(s))
	if err != nil {
		return false
	}
	return true
}

func getInitial(s string) string {
	e := []rune(s)
	return string(e[0])
}

func diffSlice(slice1 []string, slice2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

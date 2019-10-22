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

	"github.com/iancoleman/strcase"
	"github.com/shogo82148/txmanager"
)

var (
	database  = flag.String("d", "", "Name of Database")
	host      = flag.String("h", "127.0.0.1", "Host of Database")
	port      = flag.Int("P", 3306, "Database port number")
	user      = flag.String("u", "root", "Database username")
	password  = flag.String("p", "", "Database password")
	specific  = flag.String("S", "", "Import specific tables")
	separate  = flag.Bool("s", false, "Separate CSV into 2 types")
	ignore    = flag.Bool("i", false, "Ignore 1st line of CSV")
	auto      = flag.Bool("a", false, "Auto completion with file name when lack of csv columns")
	snakecase = flag.Bool("sn", false, "If csv columns is camelcase, convert to snakecase")
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
	baseAbsPath, _ := filepath.Abs(dir)

	csvs := csv.FindCsv(baseAbsPath)
	tables := createTableList(csvs, baseAbsPath)

	csvs = filterSpecificTables(csvs, *specific)
	tables = filterSpecificTables(tables, *specific)

	dbm := txmanager.NewDB(DB)

	var err error
	txmanager.Do(dbm, func(tx txmanager.Tx) error {
		txmanager.Do(tx, func(tx txmanager.Tx) error {
			for i, csvAbsPath := range csvs {
				mysql.RegisterLocalFile(csvAbsPath)

				dbColumns := db.GetColumns(DB, tables[i])
				csvColumns := csv.GetColumns(csvAbsPath)

				// ToSnakeCase
				var setColumns, sqlColumns, setQuery string
				if *snakecase {
					csvCamelColumns := csvColumns
					snakeColumns := toSnakeSlice(csvColumns)
					tmpColumns := connectEqual(snakeColumns, addPrefix(csvColumns, "@")) // [id=@id user_id=@userId]
					csvColumns = toSnakeSlice(csvColumns)
					setColumns = strings.Join(tmpColumns, ",")                                  // "id=@id,user_id=@userId"
					sqlColumns = "(" + strings.Join(addPrefix(csvCamelColumns, "@"), ",") + ")" // (@id,@userId)
					setQuery = sqlColumns + " SET " + setColumns                                // "(@id,@userId) SET id=@id,user_id=@userId"
				}

				diffColumns := diffSlice(dbColumns, csvColumns)

				baseQuery := "LOAD DATA LOCAL INFILE '" + csvAbsPath + "' INTO TABLE " + tables[i] + " FIELDS TERMINATED BY ',' "
				if *ignore {
					baseQuery += " IGNORE 1 LINES "
				}

				csvRelPath, _ := filepath.Rel(baseAbsPath, csvAbsPath)
				if len(diffColumns) == 0 {
					_, err = tx.Exec(baseQuery + setQuery)
					fmt.Println(csvRelPath, "import to", tables[i])
				} else if len(diffColumns) != 0 && *auto {
					csvFile := getFileNameWithoutExt(csvAbsPath)
					var sets string
					for i, column := range diffColumns {
						sets += column + "=" + csvFile
						if i != (len(diffColumns) - 1) {
							sets += ","
						}
					}
					if *snakecase {
						setQuery += "," + sets
					}

					_, err = tx.Exec(baseQuery + setQuery)
					fmt.Println(csvRelPath, "import to", tables[i])
				}

				if err != nil {
					fmt.Println("Failed: ", csvAbsPath, "->", tables[i])
					tx.TxRollback()
					log.Fatalf("Error: Query faild %v", err)
				}
			}
			return err
		})
		return err
	})

	fmt.Println("Complete !!")
}

func createTableList(paths []string, baseAbsPath string) (tableList []string) {
	for _, path := range paths {
		rpath, _ := filepath.Rel(baseAbsPath, path)

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

func addPrefix(srcSlice []string, p string) (destSlice []string) {
	for _, s := range srcSlice {
		destSlice = append(destSlice, p+s)
	}
	return destSlice
}

func connectEqual(aSlice, bSlice []string) (destSlice []string) {
	if len(aSlice) != len(bSlice) {
		log.Fatal("Error: miss ")
	}

	for i := 0; i < len(aSlice); i++ {
		destSlice = append(destSlice, aSlice[i]+"="+bSlice[i])
	}
	return destSlice
}

func toSnakeSlice(s []string) (snakeSlice []string) {
	for _, e := range s {
		snakeSlice = append(snakeSlice, strcase.ToSnake(e))
	}
	return snakeSlice
}

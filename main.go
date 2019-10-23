package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sanshirookazaki/csv2sql/csv"
	"github.com/sanshirookazaki/csv2sql/db"
	"github.com/sanshirookazaki/csv2sql/util"

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
	snakecase = flag.Int("sn", 0, "Convert columns into snakecase")
	dryrun    = flag.Bool("dry-run", false, "dry run")
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

	csvAbsPaths := csv.FindCsv(baseAbsPath)
	tables := createTables(csvAbsPaths, baseAbsPath)

	csvAbsPaths = util.FilterSpecific(csvAbsPaths, *specific)
	tables = util.FilterSpecific(tables, *specific)

	dbm := txmanager.NewDB(DB)

	var err error
	txmanager.Do(dbm, func(tx txmanager.Tx) error {
		txmanager.Do(tx, func(tx txmanager.Tx) error {
			for i, csvAbsPath := range csvAbsPaths {
				mysql.RegisterLocalFile(csvAbsPath)

				dbColumns := db.GetColumns(DB, tables[i])
				csvColumns := csv.GetColumns(csvAbsPath)

				// ToSnakeCase
				var setColumns, sqlColumns, setQuery string
				if *snakecase != 0 {
					csvCamelColumns := csvColumns
					snakeColumns := util.ToSnakeSlice(csvColumns, *snakecase)
					tmpColumns := util.ConnectEqual(snakeColumns, util.AddPrefix(csvColumns, "@")) // [id=@id user_id=@userId]
					csvColumns = util.ToSnakeSlice(csvColumns, *snakecase)
					setColumns = strings.Join(tmpColumns, ",")                                       // "id=@id,user_id=@userId"
					sqlColumns = "(" + strings.Join(util.AddPrefix(csvCamelColumns, "@"), ",") + ")" // (@id,@userId)
					setQuery = sqlColumns + " SET " + setColumns                                     // "(@id,@userId) SET id=@id,user_id=@userId"
				}

				diffColumns := util.DiffSlice(dbColumns, csvColumns)

				baseQuery := "LOAD DATA LOCAL INFILE '" + csvAbsPath + "' INTO TABLE " + tables[i] + " FIELDS TERMINATED BY ',' "
				if *ignore {
					baseQuery += " IGNORE 1 LINES "
				}

				csvRelPath, _ := filepath.Rel(baseAbsPath, csvAbsPath)
				if len(diffColumns) == 0 {
					_, err = tx.Exec(baseQuery + setQuery)
					fmt.Println(csvRelPath, "import to", tables[i])
				} else if len(diffColumns) != 0 && *auto {
					csvFile := util.GetFileNameWithoutExt(csvAbsPath)
					var sets string
					for i, column := range diffColumns {
						sets += column + "=" + csvFile
						if i != (len(diffColumns) - 1) {
							sets += ","
						}
					}
					if *snakecase != 0 {
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

		if *dryrun {
			tx.TxRollback()
		}
		return err
	})

	fmt.Println("Complete !!")
}

func createTables(targetAbsPaths []string, baseAbsPath string) (tables []string) {
	for _, targetPath := range targetAbsPaths {
		relPath, _ := filepath.Rel(baseAbsPath, targetPath)

		var table string
		if *separate && !util.InitialIsInt(filepath.Base(relPath)) {
			dir, file := filepath.Split(relPath)
			file = util.GetFileNameWithoutExt(file)
			table = dir + file
		} else {
			table = filepath.Dir(relPath)
		}

		if table == "." {
			table = filepath.Base(baseAbsPath)
		}
		tables = append(tables, strings.Replace(table, "/", "_", -1))
	}

	return tables
}

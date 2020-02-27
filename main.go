package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/sanshirookazaki/csv2sql/csv"
	"github.com/sanshirookazaki/csv2sql/db"
	"github.com/sanshirookazaki/csv2sql/util"

	"github.com/fatih/color"
	"github.com/shogo82148/txmanager"
)

var (
	database  = flag.String("d", "", "Name of Database")
	host      = flag.String("h", "127.0.0.1", "Host of Database")
	port      = flag.Int("P", 3306, "Database port number")
	user      = flag.String("u", "root", "Database username")
	password  = flag.String("p", "", "Database password")
	specific  = flag.String("S", "", "Import specific tables")
	newline   = flag.String("n", "\n", "Write newlines at the ends of lines")
	null      = flag.Bool("N", false, "If csv value is empty, set null")
	separate  = flag.Bool("s", false, "Separate CSV into 2 types")
	ignore    = flag.Bool("i", false, "Ignore 1st line of CSV")
	auto      = flag.Bool("a", false, "Auto completion with file name when lack of csv columns")
	snakecase = flag.Int("sn", 0, "Convert columns into snakecase")
	dryrun    = flag.Bool("dry-run", false, "dry run")
	force     = flag.Bool("f", false, "Force run: ignore error")
)

func main() {
	log.SetOutput(io.MultiWriter(os.Stdout))
	flag.Parse()
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", *user, *password, *host, *port, *database)
	DB := db.Conn(conn)
	defer DB.Close()

	if len(os.Args) < 2 {
		log.Panicf("Error: CSV path is required")
	}

	dir := os.Args[len(os.Args)-1]
	baseAbsPath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Error: CSV path doesn't exist %v", err)
	}

	csvAbsPaths := csv.FindCsv(baseAbsPath, *specific)
	targetTables := createTables(csvAbsPaths, baseAbsPath)
	tables, _ := db.GetTables(DB)

	dbm := txmanager.NewDB(DB)

	// stdout
	var importList []string
	var skipList []string
	var failList []string

	txmanager.Do(dbm, func(tx txmanager.Tx) error {
		for i := 0; i < len(csvAbsPaths); i++ {
			fy := color.New(color.FgYellow)
			if !util.Contains(tables, targetTables[i]) {
				fy.Println("Skip :table not exist", targetTables[i], csvAbsPaths[i]+"\n")
				skipList = append(skipList, targetTables[i])
				continue
			}
			if !csv.ExistData(csvAbsPaths[i]) {
				fy.Println("Skip :data not exist", targetTables[i], csvAbsPaths[i]+"\n")
				skipList = append(skipList, targetTables[i])
				continue
			}

			mysql.RegisterLocalFile(csvAbsPaths[i])

			dbColumns, _ := db.GetColumns(DB, targetTables[i])
			csvColumns, _ := csv.GetColumns(csvAbsPaths[i])

			// ToSnakeCase
			var setColumns, sqlColumns, setQuery string
			if *snakecase != 0 {
				csvCamelColumns := csvColumns
				snakeColumns := util.ToSnakeSlice(csvColumns, *snakecase)
				var tmpColumns []string
				if *null {
					tmpColumns = util.ConnectEqual(util.EncloseMark(snakeColumns, "`", "`"), util.SetNullValue(csvColumns)) // [`id`= case @id when '' then NULL else @id end `user_id`= case @userId when '' then NULL else @userId end ]
				} else {
					tmpColumns = util.ConnectEqual(util.EncloseMark(snakeColumns, "`", "`"), util.AddPrefix(csvColumns, "@")) // [`id`=@id `user_id`=@userId]
				}
				csvColumns = util.ToSnakeSlice(csvColumns, *snakecase)
				setColumns = strings.Join(tmpColumns, ",")                                       // "`id`=@id,`user_id`=@userId"
				sqlColumns = "(" + strings.Join(util.AddPrefix(csvCamelColumns, "@"), ",") + ")" // (@id,@userId)
				setQuery = sqlColumns + " SET " + setColumns                                     // "(@id,@userId) SET `id`=@id,`user_id`=@userId"
			}

			diffColumns := util.DiffSlice(dbColumns, csvColumns)
			diffColumns = util.RemoveElements(diffColumns, []string{"created_at", "updated_at"})

			baseQuery := "LOAD DATA LOCAL INFILE '" + csvAbsPaths[i] + "' INTO TABLE " + targetTables[i] + " FIELDS TERMINATED BY ',' LINES TERMINATED BY " + "'" + *newline + "'"
			if *ignore {
				baseQuery += " IGNORE 1 LINES "
			}

			csvRelPath, err := filepath.Rel(baseAbsPath, csvAbsPaths[i])
			if err != nil {
				log.Fatalf("Error: Can't create CsvRelPath %v", err)
			}
			var query string
			if len(diffColumns) == 0 {
				query = baseQuery + setQuery
			} else if len(diffColumns) != 0 && *auto {
				csvFile := util.GetFileNameWithoutExt(csvAbsPaths[i])
				var sets string
				for i, column := range diffColumns {
					sets += "`" + column + "`='" + csvFile + "'"
					if i != (len(diffColumns) - 1) {
						sets += ","
					}
				}
				if *snakecase != 0 {
					setQuery += "," + sets
				}

				query = baseQuery + setQuery
			}

			err = db.TxExecQuery(tx, query, *force)
			if err != nil {
				fr := color.New(color.FgRed)
				fr.Println("Failed: ", csvRelPath, "->", targetTables[i], "\n"+"Rollback")
				failList = append(failList, targetTables[i])
				tx.TxRollback()
				log.Fatalf("Error: Query faild %v", err)
			}
			fg := color.New(color.FgGreen)
			fg.Println(csvRelPath, "import to", targetTables[i]+"\n")
			importList = append(importList, targetTables[i])
		}

		if *dryrun {
			tx.TxRollback()
			fmt.Println("Dry Run !")
		}
		return err
	})
	if len(importList) > 0 {
		i := color.New(color.FgGreen)
		i.Println("Success :", importList)
	}
	if len(skipList) > 0 {
		s := color.New(color.FgYellow)
		s.Println("Skip :", skipList)
	}
	if len(failList) > 0 {
		f := color.New(color.FgRed)
		f.Println("Failed :", failList)
	}
	fc := color.New(color.FgCyan)
	fc.Println("Complete !!")
}

func createTables(targetAbsPaths []string, baseAbsPath string) (tables []string) {
	for _, targetPath := range targetAbsPaths {
		relPath, err := filepath.Rel(baseAbsPath, targetPath)
		if err != nil {
			log.Fatalf("Error: Can't create RelPath %v", err)
		}

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

package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func GetColumns(csvPath string) (columns []string) {
	file, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	columns, err = reader.Read()
	if err != nil {
		log.Panicf("Error: Can't read csv columns")
	}

	return columns
}

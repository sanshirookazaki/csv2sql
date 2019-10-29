package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func GetColumns(csvPath string) (columns []string, err error) {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Error: Can't Open file %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	columns, err = reader.Read()
	if err != nil {
		log.Printf("Error: Can't read csv columns %v", err)
	}

	return columns, err
}

func ExistData(csvPath string) bool {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Error: Can't Open file %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for i := 0; i < 2; i++ {
		_, err = reader.Read()
	}

	if err != nil {
		return false
	}
	return true
}

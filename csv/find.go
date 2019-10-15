package csv

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func FindCsv(dir string) (paths []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panicf("Error: Path doesn't exist")
	}

	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, FindCsv(filepath.Join(dir, file.Name()))...)
			continue
		}

		if strings.HasSuffix(file.Name(), ".csv") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}

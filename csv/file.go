package csv

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func FindCsv(dir, specific string) (paths []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panicf("Error: Path doesn't exist")
	}

	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, FindCsv(filepath.Join(dir, file.Name()), specific)...)
			continue
		}

		if strings.HasSuffix(file.Name(), ".csv") && strings.Contains(dir+file.Name(), specific) {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths
}

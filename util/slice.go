package util

import (
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

func FilterSpecific(list []string, specific string) (result []string) {
	for _, e := range list {
		if strings.Contains(e, specific) {
			result = append(result, e)
		}
	}
	return result
}

func DiffSlice(slice1 []string, slice2 []string) []string {
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

func AddPrefix(srcSlice []string, p string) (destSlice []string) {
	for _, s := range srcSlice {
		destSlice = append(destSlice, p+s)
	}
	return destSlice
}

func ConnectEqual(aSlice, bSlice []string) (destSlice []string) {
	if len(aSlice) != len(bSlice) {
		log.Fatal("Error: miss ")
	}

	for i := 0; i < len(aSlice); i++ {
		destSlice = append(destSlice, aSlice[i]+"="+bSlice[i])
	}
	return destSlice
}

func ToSnakeSlice(s []string) (snakeSlice []string) {
	for _, e := range s {
		snakeSlice = append(snakeSlice, strcase.ToSnake(e))
	}
	return snakeSlice
}
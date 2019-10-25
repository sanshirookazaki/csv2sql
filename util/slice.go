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
		log.Fatal("Error: Miss ConnectEqual, Different Length Slice")
	}

	for i := 0; i < len(aSlice); i++ {
		destSlice = append(destSlice, aSlice[i]+"="+bSlice[i])
	}
	return destSlice
}

func ToSnakeSlice(s []string, i int) (snakeSlice []string) {
	for _, e := range s {
		if i == 1 {
			e = strcase.ToSnake(e)
		} else if i == 2 {
			e = ToSnake(e)
		}
		snakeSlice = append(snakeSlice, e)
	}
	return snakeSlice
}

func RemoveElements(target, remove []string) (dest []string) {
	for i := 0; i < len(target); i++ {
		tar := target[i]
		for _, rem := range remove {
			if tar == rem {
				dest = append(target[:i], target[i+1:]...)
				i--
				break
			}
		}
	}
	return dest
}

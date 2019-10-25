package util

import (
	"strconv"
	"unicode"
)

func InitialIsInt(s string) bool {
	_, err := strconv.Atoi(GetInitial(s))
	if err != nil {
		return false
	}
	return true
}

func GetInitial(s string) string {
	e := []rune(s)
	return string(e[0])
}

func ToSnake(s string) string {
	runes := []rune(s)
	length := len(runes)

	var result []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(runes[i]))
	}

	return string(result)
}

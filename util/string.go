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

func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

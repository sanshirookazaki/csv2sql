package util

import "strconv"

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

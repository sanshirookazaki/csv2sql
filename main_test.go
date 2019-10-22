package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateTables(t *testing.T) {
	targetAbsPaths := []string{"./path/user/1.csv", "./path/user/task/1.csv"}
	baseAbsPath := "./path"
	tables := createTables(targetAbsPaths, baseAbsPath)
	e := []string{"user", "user_task"}
	fmt.Println(e, tables)
	if !reflect.DeepEqual(tables, e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass createTableList func")
}

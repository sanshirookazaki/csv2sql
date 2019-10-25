package main

import (
	"reflect"
	"testing"
)

func TestCreateTables(t *testing.T) {
	targetAbsPaths := []string{"./path/user/1.csv", "./path/user/task/1.csv"}
	baseAbsPath := "./path"
	tables := createTables(targetAbsPaths, baseAbsPath)
	e := []string{"user", "user_task"}
	if !reflect.DeepEqual(tables, e) {
		t.Fatalf("error: %v", e)
	}

	*separate = true
	targetAbsPaths2 := []string{"./path/detail.csv", "./path/user/detail.csv", "./path/user/task/detail.csv"}
	baseAbsPath2 := "./path"
	tables2 := createTables(targetAbsPaths2, baseAbsPath2)
	e2 := []string{"detail", "user_detail", "user_task_detail"}
	if !reflect.DeepEqual(tables2, e2) {
		t.Fatalf("error: %v", e2)
	}

	t.Log("Pass createTables func")
}

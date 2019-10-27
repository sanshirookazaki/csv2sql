package csv

import (
	"reflect"
	"testing"
)

func TestFindCsv(t *testing.T) {
	dir := "../examples"
	specific := "user/detail"
	e := []string{"user/detail.csv"}
	if reflect.DeepEqual(FindCsv(dir, specific), e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass FindCsv func")
}

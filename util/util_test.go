package util

import (
	"reflect"
	"testing"
)

func TestGetInitial(t *testing.T) {
	e := GetInitial("abc")
	if !(e == "a") {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass GetInitial func")
}

func TestInitialIsInt(t *testing.T) {
	e1 := InitialIsInt("123abc")
	if !e1 {
		t.Fatalf("error: %v", e1)
	}
	e2 := InitialIsInt("abc")
	if e2 {
		t.Fatalf("error: %v", e2)
	}
	t.Log("Pass InitialIsInt func")
}

func TestFilterSpecific(t *testing.T) {
	list := []string{"userId", "userTask", "time"}
	specific := "user"
	result := FilterSpecific(list, specific)
	e := []string{"userId", "userTask"}
	if !reflect.DeepEqual(result, e) {
		{
			t.Fatalf("error: %v", e)
		}
		t.Log("Pass FilterSpecific func")
	}
}

func TestDiffSlice(t *testing.T) {
	slice1 := []string{"userId", "userTask", "time"}
	slice2 := []string{"userId", "time"}
	result := DiffSlice(slice1, slice2)
	e := []string{"userTask"}
	if !reflect.DeepEqual(result, e) {
		{
			t.Fatalf("error: %v", e)
		}
		t.Log("Pass DiffSlice func")
	}
}

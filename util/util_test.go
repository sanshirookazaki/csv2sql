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

func TestGetFileNameWithoutExt(t *testing.T) {
	path := "/foo/bar/test1.csv"
	e := "test1"
	if !(GetFileNameWithoutExt(path) == e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass GetFileNameWithoutExt func")
}

func TestToSnake(t *testing.T) {
	s := "aaaBBB123CCC"
	e := "aaa_bbb123_ccc"
	if ToSnake(s) != e {
		t.Fatalf("error: %v", ToSnake(s))
	}
	t.Log("Pass ToSnake fund")
}

func TestAddPrefix(t *testing.T) {
	s := []string{"test", "abc"}
	p := "@"
	e := []string{"@test", "@abc"}
	if !reflect.DeepEqual(AddPrefix(s, p), e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass AddPrefix func")
}

func TestConnectEqual(t *testing.T) {
	slice1 := []string{"@test1", "@test2"}
	slice2 := []string{"test1", "test2"}
	e := []string{"@test1=test1", "@test2=test2"}
	if !reflect.DeepEqual(ConnectEqual(slice1, slice2), e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass ConnectEqual func")
}

func TestToSnakeSlice(t *testing.T) {
	s1 := []string{"aaaBBB123CCC"}
	e1 := []string{"aaa_bbb_123_ccc"}
	e2 := []string{"aaa_bbb123_ccc"}
	if !reflect.DeepEqual(ToSnakeSlice(s1, 1), e1) {
		t.Fatalf("error1: %v", e1)
	}
	if !reflect.DeepEqual(ToSnakeSlice(s1, 2), e2) {
		t.Fatalf("error2: %v", e2)
	}
	t.Log("Pass ToSnakeSlice func")
}

func TestRemoveElements(t *testing.T) {
	s := []string{"test", "abc", "one", "123"}
	r := []string{"one", "test123", "two"}
	e := []string{"test", "abc"}
	if !reflect.DeepEqual(RemoveElements(s, r), e) {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass RemoveElemnts func")
}

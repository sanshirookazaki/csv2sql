package util

import "testing"

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

package main

import "testing"

func TestGetInitial(t *testing.T) {
	e := getInitial("abc")
	if !(e == "a") {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass getInitial func")
}

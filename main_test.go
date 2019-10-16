package main

import "testing"

func TestgetInitial(t *testing.T) {
	e := getInitial("abc")
	if !(e == "a") {
		t.Fatalf("error: %v", e)
	}
	t.Log("Pass getInitial func")
}

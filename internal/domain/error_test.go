package domain

import "testing"

func TestError_ErrorMethod(t *testing.T) {
	e := Error{Code: "X", Message: "hello"}
	if e.Error() != "hello" {
		t.Fatal("want message")
	}
}

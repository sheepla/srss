package ui

import (
	"testing"
)

func TestHumanizeTime(t *testing.T) {
	have := humanizeTime(nil)
	want := ""
	if have != want {
		t.Error("e!")
	}
}

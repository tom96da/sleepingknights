package main

import (
	"os"
	"testing"
)

func TestRunMainHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"slk", "--help"}

	code := runMain()
	if code != 0 {
		t.Fatalf("unexpected exit code: got=%d want=0", code)
	}
}

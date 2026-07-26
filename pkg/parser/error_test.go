package parser

import "testing"

func TestParseErrorStringer(t *testing.T) {
	err := ParseError{Line: 3, Column: 4, Message: "oops"}
	if err.Error() != "3:4: oops" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
}

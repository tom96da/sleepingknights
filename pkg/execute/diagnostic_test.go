package execute

import "testing"

func TestBuildSourceDiagnostic(t *testing.T) {
	src := "let x = 1\nprint(x)\n"
	got := buildSourceDiagnostic("examples/main.slk", src, 2, 7, "undefined variable 'y'")

	want := "examples/main.slk:2:7: error: undefined variable 'y'\n" +
		"  | print(x)\n" +
		"  |       ^"

	if got != want {
		t.Fatalf("diagnostic mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestBuildSourceDiagnosticBounds(t *testing.T) {
	src := "abc\n"
	got := buildSourceDiagnostic("a.slk", src, 0, 0, "msg")
	want := "a.slk:1:1: error: msg\n  | abc\n  | ^"
	if got != want {
		t.Fatalf("diagnostic bounds mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

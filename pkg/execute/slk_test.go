package execute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDisplayHelpers(t *testing.T) {
	if displayUsage() != ExitStatusSuccess {
		t.Fatal("displayUsage status")
	}
	if displayCommandHelp("run") != ExitStatusSuccess {
		t.Fatal("displayCommandHelp(run) status")
	}
	if displayCommandHelp("unknown") != ExitStatusInvalidCommandLineArgs {
		t.Fatal("displayCommandHelp(unknown) status")
	}
	if displayVersion() != ExitStatusSuccess {
		t.Fatal("displayVersion status")
	}
	if displayArgsError("x") != ExitStatusInvalidCommandLineArgs {
		t.Fatal("displayArgsError status")
	}
	if startInteractiveMode() != ExitStatusSuccess {
		t.Fatal("startInteractiveMode status")
	}
}

func TestIsScriptFile(t *testing.T) {
	if !isScriptFile("a.slk") {
		t.Fatal("expected .slk true")
	}
	if isScriptFile("a.txt") {
		t.Fatal("expected non-.slk false")
	}
	if isScriptFile("abc") {
		t.Fatal("expected short false")
	}
}

func TestCommandLineBranches(t *testing.T) {
	dir := t.TempDir()
	valid := filepath.Join(dir, "ok.slk")
	if err := os.WriteFile(valid, []byte("let x = 1\n"), 0o644); err != nil {
		t.Fatalf("write valid: %v", err)
	}
	lexErr := filepath.Join(dir, "lex.slk")
	if err := os.WriteFile(lexErr, []byte("@\n"), 0o644); err != nil {
		t.Fatalf("write lex err: %v", err)
	}
	parseErr := filepath.Join(dir, "parse.slk")
	if err := os.WriteFile(parseErr, []byte("let x =\n"), 0o644); err != nil {
		t.Fatalf("write parse err: %v", err)
	}

	cases := []struct {
		name string
		args []string
		want ExitStatus
	}{
		{"interactive", []string{}, ExitStatusSuccess},
		{"help", []string{"--help"}, ExitStatusSuccess},
		{"version", []string{"--version"}, ExitStatusSuccess},
		{"run-help", []string{"run", "--help"}, ExitStatusSuccess},
		{"unknown-help", []string{"wat", "--help"}, ExitStatusInvalidCommandLineArgs},
		{"run-missing", []string{"run"}, ExitStatusInvalidCommandLineArgs},
		{"run-not-implemented", []string{"run", valid}, ExitStatusNotImplemented},
		{"direct-valid", []string{valid}, ExitStatusSuccess},
		{"direct-lex-error", []string{lexErr}, ExitStatusCriticalException},
		{"direct-parse-error", []string{parseErr}, ExitStatusCriticalException},
		{"unknown", []string{"wat"}, ExitStatusInvalidCommandLineArgs},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CommandLine(tc.args)
			if got != tc.want {
				t.Fatalf("status mismatch: got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestExecuteScriptReadError(t *testing.T) {
	got := executeScript(filepath.Join(t.TempDir(), "missing.slk"))
	if got != ExitStatusCriticalException {
		t.Fatalf("status mismatch: got=%v want=%v", got, ExitStatusCriticalException)
	}
}

func TestCompileAndExecuteStatus(t *testing.T) {
	if compileAndExecute("x.slk") != ExitStatusNotImplemented {
		t.Fatal("compileAndExecute status")
	}
}

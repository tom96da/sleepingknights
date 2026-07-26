package execute

import (
	"fmt"

	"github.com/tom96da/sleepingknights/pkg"
)

func displayUsage() ExitStatus {
	fmt.Println("Usage: slk [options] [script.slk] | slk run <script.slk>")
	fmt.Println()
	fmt.Println("  slk                  Start REPL (interactive mode)")
	fmt.Println("  slk <script.slk>     Run script in interpreter mode")
	fmt.Println()
	fmt.Println("Command stubs:")
	fmt.Println("  run <script.slk>     Compile and run script (planned for Phase 6)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("  -v, --version        Show version information")

	return ExitStatusSuccess
}

func displayCommandHelp(command string) ExitStatus {
	switch command {
	case "run":
		fmt.Println("Usage: slk run <script.slk>")
		fmt.Println()
		fmt.Println("Compile and run a script file.")
		fmt.Println("Status: planned for Phase 6 (not implemented yet)")
	default:
		return displayArgsError(fmt.Sprintf("Unknown command: %s", command))
	}

	return ExitStatusSuccess
}

func displayVersion() ExitStatus {
	fmt.Printf("SleepingKnights  %s (%s)\n", pkg.Version, pkg.Hash)
	return ExitStatusSuccess
}

func displayArgsError(message string) ExitStatus {
	fmt.Printf("[Error] %s\n", message)
	fmt.Println()
	displayUsage()
	return ExitStatusInvalidCommandLineArgs
}

func startInteractiveMode() ExitStatus {
	displayVersion()
	return ExitStatusSuccess
}

func CommandLine(commandLineArgs []string) ExitStatus {
	if len(commandLineArgs) == 0 {
		return startInteractiveMode()
	}

	if len(commandLineArgs) >= 2 && (commandLineArgs[1] == "-h" || commandLineArgs[1] == "--help") {
		return displayCommandHelp(commandLineArgs[0])
	}

	switch commandLineArgs[0] {
	case "-h", "--help":
		return displayUsage()
	case "-v", "--version":
		return displayVersion()
	case "run":
		if len(commandLineArgs) < 2 {
			return displayArgsError("Missing script path")
		}
		return compileAndExecute(commandLineArgs[1])
	default:
		if isScriptFile(commandLineArgs[0]) {
			return executeScript(commandLineArgs[0])
		}
		return displayArgsError(fmt.Sprintf("Unknown command: %s", commandLineArgs[0]))
	}
}

func isScriptFile(path string) bool {
	return len(path) > 4 && path[len(path)-4:] == ".slk"
}

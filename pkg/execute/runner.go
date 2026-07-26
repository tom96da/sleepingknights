package execute

import (
	"fmt"
	"os"

	"github.com/tom96da/sleepingknights/pkg/lexer"
	"github.com/tom96da/sleepingknights/pkg/parser"
	"github.com/tom96da/sleepingknights/pkg/token"
)

func executeScript(scriptPath string) ExitStatus {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Printf("[Error] Failed to read script: %v\n", err)
		return ExitStatusCriticalException
	}

	l := lexer.New(string(content))
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		if tok.Type == token.ILLEGAL {
			fmt.Println(buildSourceDiagnostic(scriptPath, string(content), tok.Line, tok.Column, tok.Literal))
			return ExitStatusCriticalException
		}
	}

	p := parser.New(string(content))
	_, parseErrs := p.ParseProgram()
	if len(parseErrs) > 0 {
		for _, parseErr := range parseErrs {
			fmt.Println(buildSourceDiagnostic(scriptPath, string(content), parseErr.Line, parseErr.Column, parseErr.Message))
		}
		return ExitStatusCriticalException
	}

	return ExitStatusSuccess
}

func compileAndExecute(_ string) ExitStatus {
	return ExitStatusNotImplemented
}

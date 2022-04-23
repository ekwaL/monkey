package runner

import (
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

func RunFile(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Can not read file %q : %s", name, err.Error()))
		os.Exit(64)
	}
	runProgram(string(data))
}

func runProgram(source string) {
	env := object.NewEnvironment()
	l := lexer.New(string(source))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParseErrors(os.Stdout, p.Errors())
	}

	eval.Eval(program, env)
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Errors found while parsing:\n")
	for _, e := range errors {
		io.WriteString(out, "\t"+e+"\n")
	}
}

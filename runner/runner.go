package runner

import (
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/resolver"
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
		printParseErrors(os.Stderr, p.Errors())
		os.Exit(65)
	}

	r := resolver.New()
	r.Resolve(program)

	if len(r.Errors()) != 0 {
		printParseErrors(os.Stderr, r.Errors())
		os.Exit(65)
	}

	eval.Locals = r.Locals()

	evalResult := eval.Eval(program, env)

	if evalResult != nil && evalResult.Type() == object.ERROR_OBJ {
		io.WriteString(os.Stderr, evalResult.Inspect()+"\n")
		os.Exit(70)
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Errors found while parsing:\n")
	for _, e := range errors {
		io.WriteString(out, "\t"+e+"\n")
	}
}

package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const ReplWelcomeMessage = `
This is a Monkey REPL.
More info on usage later.
`
const usageInfo = `
Usage:
monkey [script]  to run script
monkey           to run REPL
`

type REPL struct {
	in  io.Reader
	out io.Writer
}

func New(in io.Reader, out io.Writer) *REPL {
	return &REPL{in: in, out: out}
}

func (r *REPL) Start() {
	r.printWelcomeMessage()

	scanner := bufio.NewScanner(r.in)
	env := object.NewEnvironment()

	lineNumber := 0
	for {
		io.WriteString(r.out, prompt(lineNumber))
		lineNumber += 1
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(r.out, p.Errors())
		}
		evalResult := eval.Eval(program, env)

		if evalResult != nil {
			io.WriteString(r.out, evalResult.Inspect())
			io.WriteString(r.out, "\n")
		}
	}
}

func (r *REPL) printWelcomeMessage() {
	io.WriteString(r.out, ReplWelcomeMessage)
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Errors found while parsing:\n")
	for _, e := range errors {
		io.WriteString(out, "\t"+e+"\n")
	}
}

func prompt(lineNumber int) string {
	return fmt.Sprintf("monkey:%03d>> ", lineNumber)
}

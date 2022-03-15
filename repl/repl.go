package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			io.WriteString(r.out, fmt.Sprintf("%+v\n", tok))
		}
	}
}

func (r *REPL) printWelcomeMessage() {
	io.WriteString(r.out, ReplWelcomeMessage)
}

func prompt(lineNumber int) string {
	return fmt.Sprintf("monkey:%03d>> ", lineNumber)
}

package main

import (
	"monkey/repl"
	"os"
)

func main() {
	r := repl.New(os.Stdin, os.Stdout)
	r.Start()
}

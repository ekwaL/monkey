package main

import (
	"fmt"
	"monkey/repl"
	"monkey/runner"
	"os"
)

const usageInfo = `
Usage:
monkey [script]  to run script
monkey           to run REPL`

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		r := repl.New(os.Stdin, os.Stdout)
		r.Start()
	case 1:
		runner.RunFile(args[0])
	default:
		fmt.Println(usageInfo)
		os.Exit(64)
	}
}

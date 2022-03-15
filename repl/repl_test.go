package repl_test

import (
	"bytes"
	"monkey/repl"
	"strings"
	"testing"
)

func TestRepl(t *testing.T) {
	t.Run("it prints welcome message", func(t *testing.T) {
		r, out := createReplWithReader("")
		r.Start()

		if !strings.HasPrefix(out.String(), repl.ReplWelcomeMessage) {
			t.Errorf("no repl welcome message printed. expect %q to start with %q",
				out.String(), repl.ReplWelcomeMessage)
		}
	})

	t.Run("it prints prompt", func(t *testing.T) {
		r, out := createReplWithReader("")
		r.Start()

		prompt := "monkey:000>> "
		if !strings.HasSuffix(out.String(), prompt) {
			t.Errorf("no prompt printed. expect %q to end with %q",
				out.String(), prompt)
		}
	})
}

func createReplWithReader(input string) (*repl.REPL, *bytes.Buffer) {
	reader := strings.NewReader(input)
	out := new(bytes.Buffer)
	r := repl.New(reader, out)

	return r, out
}

package repl_test

import (
	"bufio"
	"monkey/repl"
	"monkey/token"
	"strings"
	"testing"
)

func TestReplIntegration(t *testing.T) {
	t.Run("it accepts and tokenizes input", func(t *testing.T) {
		tt := []struct {
			input  string
			tokens []token.Token
		}{
			{"let five = 5 ;", []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
			}},
			{"10 == 10;", []token.Token{
				{Type: token.INT, Literal: "10"},
				{Type: token.EQUAL_EQUAL, Literal: "=="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
			}},
			{"10 <= 11;", []token.Token{
				{Type: token.INT, Literal: "10"},
				{Type: token.LESS_EQUAL, Literal: "<="},
				{Type: token.INT, Literal: "11"},
				{Type: token.SEMICOLON, Literal: ";"},
			}},
		}

		for _, tc := range tt {
			r, out := createReplWithReader(tc.input)
			r.Start()

			scanner := bufio.NewScanner(out)

			got := ""
			for scanner.Scan() {
				got += scanner.Text() + "\n"
			}

			got = strings.TrimPrefix(got, repl.ReplWelcomeMessage+"monkey:000>> ")
			got = strings.TrimSuffix(got, "monkey:001>> \n")

			want := ""
			for _, tok := range tc.tokens {
				want += tok.String() + "\n"
			}

			if got != want {
				t.Errorf("wrong repl output. want %q, got %q", want, got)
			}
		}
	})
}

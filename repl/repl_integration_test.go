package repl_test

import (
	"bufio"
	"monkey/repl"
	"strings"
	"testing"
)

func TestReplIntegration(t *testing.T) {
	t.Run("it accepts and tokenizes input", func(t *testing.T) {
		tt := []struct {
			input string
			want  string
		}{
			{"let  five  =  5 ;", "let five = 5;"},
			{"10 == 10;", "(10 == 10);"},
			{"fn (i, j) { i + j }", "fn(i, j) { (i + j); };"},
			{"if( x > y ) true", "if (x > y) true;;"},
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
			got = strings.Trim(got, "\n")

			if got != tc.want {
				t.Errorf("wrong repl output. want %q, got %q", tc.want, got)
			}
		}
	})
}

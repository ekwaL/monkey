package repl_test

import (
	"bufio"
	"monkey/repl"
	"strings"
	"testing"
)

func TestReplIntegration(t *testing.T) {
	tt := []struct {
		input string
		want  string
	}{
		{"12345;", "12345"},
		{"true; false;", "false"},
		{"false; 12345; true;", "true"},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
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
		})
	}
}

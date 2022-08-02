package code_test

import (
	"monkey/code"
	"strings"
	"testing"
)

func TestMake(t *testing.T) {
	tt := []struct {
		op       code.OpCode
		operands []int
		want     []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
	}

	for _, tc := range tt {
		instruction := code.Make(tc.op, tc.operands...)

		if len(instruction) != len(tc.want) {
			t.Errorf("Instruction has wrong length, want %d, got %d.", len(tc.want), len(instruction))
		}

		i := 0
		for ; i < len(tc.want); i++ {
			b := tc.want[i]

			if i >= len(instruction) {
				t.Errorf("Want byte %d at pos %d, got nothing.", b, i)
				continue
			}

			if instruction[i] != b {
				t.Errorf("Wrong byte at pos %d, want %d, got %d.", i, b, instruction[i])
			}
		}

		for ; i < len(instruction); i++ {
			t.Errorf("Got byte %d at pos %d, want nothing.", instruction[i], i)
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 1),
		code.Make(code.OpConstant, 10),
		code.Make(code.OpConstant, 65535),
	}

	want := `0000 OpConstant 1
0003 OpConstant 10
0006 OpConstant 65535
`

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	got := concatted.String()
	if got != want {
		t.Errorf("Wrong instructions.String(), want %q, got %q.", want, got)

		gotLines := strings.Split(got, "\n")
		wantLines := strings.Split(want, "\n")
		i := 0

		for ; i < len(gotLines); i++ {
			if i >= len(wantLines) {
				t.Errorf("%02d line want nothing", i)
				t.Errorf("%02d line  got: %q", i, gotLines[i])
				continue
			}
			if gotLines[i] != wantLines[i] {
				t.Errorf("%02d line want: %q", i, wantLines[i])
				t.Errorf("%02d line  got: %q", i, gotLines[i])
			}
		}

		for ; i < len(wantLines); i++ {
			t.Errorf("%02d line want: %q", i, wantLines[i])
			t.Errorf("%02d line got nothing", i)
		}
	}
}

func TestReadOperands(t *testing.T) {
	tt := []struct {
		op        code.OpCode
		operands  []int
		bytesRead int
	}{
		{code.OpConstant, []int{65535}, 2},
	}

	for _, tc := range tt {
		instruction := code.Make(tc.op, tc.operands...)

		def, err := code.Lookup(byte(tc.op))
		if err != nil {
			t.Fatalf("Definition is not found when it should: %q.", err)
		}

		operandsRead, n := code.ReadOperands(def, instruction[1:])
		if n != tc.bytesRead {
			t.Fatalf("Wrong amount of byted read, want %d, got %d.", tc.bytesRead, n)
		}

		for i, want := range tc.operands {
			if operandsRead[i] != want {
				t.Errorf("Wrong operand, want %d, got %d", want, operandsRead[i])
			}
		}

	}
}

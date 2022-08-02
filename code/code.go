package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Instructions []byte

type OpCode byte

const (
	_ OpCode = iota
	OpConstant
)

func (ins Instructions) String() string {
	var out bytes.Buffer

	offset := 0
	for offset < len(ins) {
		op := ins[offset]

		def, err := Lookup(op)
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			offset++
			continue
		}

		byteOperands := ins[offset+1:]
		operands, bytesRead := ReadOperands(def, byteOperands)

		strOperands := []string{}
		for _, operand := range operands {
			strOperands = append(strOperands, strconv.Itoa(operand))
		}

		out.WriteString(fmt.Sprintf("%04d %s %s\n", offset, def.Name, strings.Join(strOperands, ", ")))

		offset += bytesRead + 1 // + OpCode byte
	}

	return out.String()
}

// func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
// 	operandCount := len(def.OperandWidths)

// 	if len(operands) != operandCount {
// 		return fmt.Sprintf("ERROR: operand cound %d does not match defined %d\n",
// 			len(operands), operandCount)
// 	}

// 	swtich operandCount {
// 	case 1:
// 		return fmt.Sprintf("%s %d", def.Name, operands[0])
// 	}

// 	return fmt.Sprintf("ERROR: unhandled operand count for %s\n", def.Name)
// }

func ReadOperands(def *Definition, instruction Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(instruction[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func Make(op OpCode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	operandsLen := len(def.OperandWidths)
	if len(operands) != operandsLen {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // one operand with 2 byte width (uint16)
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d is undefined", op)
	}
	return def, nil
}

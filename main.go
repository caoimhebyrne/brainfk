package main

import (
	"bufio"
	"fmt"
	"os"
)

type Context struct {
	cells       []byte
	dataPointer uint16
	memorySize  uint16
}

func NewContext() *Context {
	return &Context{cells: make([]byte, 30000), memorySize: 30000, dataPointer: 0}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: brainfk <file>")
		os.Exit(-1)
		return
	}

	fileName := os.Args[1]
	lexer, err := NewLexer(fileName)
	if err != nil {
		fmt.Println("ERROR: Failed to create a lexer for `", fileName, "`.")
		return
	}

	fmt.Println("INFO: Lexer initialized with", len(lexer.contents), "bytes")

	instructions, err := ParseInstructions(lexer)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	context := NewContext()
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]

		switch instruction.instructionType {
		case Instruction_Inc:
			context.cells[context.dataPointer] += byte(instruction.value)

		case Instruction_Dec:
			context.cells[context.dataPointer] -= byte(instruction.value)

		case Instruction_Add:
			amount := uint16(instruction.value)
			if context.dataPointer == (context.memorySize - amount) {
				context.dataPointer = 0
			} else {
				context.dataPointer += amount
			}

		case Instruction_Sub:
			amount := uint16(instruction.value)
			if context.dataPointer == 0 {
				context.dataPointer = context.memorySize - amount
			} else {
				context.dataPointer -= amount
			}

		case Instruction_Output:
			fmt.Printf("%c", context.cells[context.dataPointer])

		case Instruction_Input:
			input, _ := reader.ReadByte()
			context.cells[context.dataPointer] = input

		case Instruction_JumpIfZero:
			if context.cells[context.dataPointer] == 0 {
				i = instruction.value
			}

		case Instruction_JumpIfNonZero:
			if context.cells[context.dataPointer] != 0 {
				i = instruction.value
			}
		}
	}
}

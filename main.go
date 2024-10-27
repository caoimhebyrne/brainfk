package main

import (
	"fmt"
	"os"
	"os/exec"
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

	transpiler := NewTranspiler(instructions)
	output := transpiler.Transpile()

	err = os.WriteFile("./build/output.S", []byte(output), 0644)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	result, err := exec.Command("/usr/bin/cc", "./build/output.S", "-o", "./build/output").Output()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	if len(result) != 0 {
		fmt.Println("INFO:", string(result[:]))
	}

	fmt.Println("INFO: Compiled to ./build/output (./build/output.S)")

	/*
		context := NewContext()
		reader := bufio.NewReader(os.Stdin)

		for i := 0; i < len(instructions); i++ {
			instruction := instructions[i]

			switch instruction.instructionType {
			case InstructionIncrement:
				context.cells[context.dataPointer] += byte(instruction.value)

			case InstructionDecrement:
				context.cells[context.dataPointer] -= byte(instruction.value)

			case InstructionAdd:
				amount := uint16(instruction.value)
				if context.dataPointer == (context.memorySize - amount) {
					context.dataPointer = 0
				} else {
					context.dataPointer += amount
				}

			case InstructionSubtract:
				amount := uint16(instruction.value)
				if context.dataPointer == 0 {
					context.dataPointer = context.memorySize - amount
				} else {
					context.dataPointer -= amount
				}

			case InstructionOutput:
				fmt.Printf("%c", context.cells[context.dataPointer])

			case InstructionInput:
				input, _ := reader.ReadByte()
				context.cells[context.dataPointer] = input

			case InstructionJumpIfZero:
				if context.cells[context.dataPointer] == 0 {
					i = instruction.value
				}

			case InstructionJumpIfNonZero:
				if context.cells[context.dataPointer] != 0 {
					i = instruction.value
				}
			}
		}
	*/
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	path := os.Args[1]
	lexer, err := NewLexer(path)
	if err != nil {
		fmt.Println("ERROR: Failed to create a lexer for `", path, "`.")
		return
	}

	fmt.Println("INFO: Lexer initialized with", len(lexer.contents), "bytes")

	instructions, err := ParseInstructions(lexer)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("INFO: Parsed", len(instructions), "instructions")

	transpiler := NewTranspiler(instructions)
	output := transpiler.Transpile()

	outputDirectory := "./build/"

	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	assemblyFile := outputDirectory + fileName + ".s"
	binaryFile := outputDirectory + fileName

	err = os.WriteFile(assemblyFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	result, err := exec.Command("/usr/bin/as", assemblyFile, "-o", binaryFile).CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: Failed to compile assembly!")
		fmt.Println(string(result[:]))
		return
	}

	fmt.Println("INFO: Compiled to", binaryFile, "("+assemblyFile+")")

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

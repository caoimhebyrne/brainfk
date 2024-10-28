package main

import (
	"fmt"
	"strings"
)

type Transpiler struct {
	instructions []Instruction
}

func NewTranspiler(instructions []Instruction) Transpiler {
	return Transpiler{instructions}
}

func (c *Transpiler) Transpile() string {
	var sb strings.Builder

	sb.WriteString(".global _main\n")
	sb.WriteString(".align 4\n")

	sb.WriteString("_main:\n")
	sb.WriteString("    adrp X1, CELLS@PAGE\n")
	sb.WriteString("    add X1, X1, CELLS@PAGEOFF\n")

	for index, instruction := range c.instructions {
		sb.WriteString(fmt.Sprintf("    // Instruction %d - %s\n", index, instruction.instructionType.String()))
		switch instruction.instructionType {
		case InstructionIncrement:
			sb.WriteString("    ldrb W0, [X1] // Load the data at X1 into W0\n")
			sb.WriteString(fmt.Sprintf("    add W0, W0, #%d\n", instruction.value))
			sb.WriteString("    strb W0, [X1] // Store the value in W0 to the memory at X1\n")

		case InstructionDecrement:
			// Load the data at X1 into X2
			// Subtract a certain value from X2
			// Store the data from X2 into the memory at X1
			sb.WriteString("    ldrb W0, [X1] // Load the data at X1 into W0\n")
			sb.WriteString(fmt.Sprintf("    sub W0, W0, #%d\n", instruction.value))
			sb.WriteString("    strb W0, [X1] // Store the value in W0 to the memory at X1\n")

		case InstructionAdd:
			// Increment X1 by a certain value.
			sb.WriteString(fmt.Sprintf("    add X1, X1, #%d\n", instruction.value))

		case InstructionSubtract:
			// Decrement X1 by a certain value.
			sb.WriteString(fmt.Sprintf("    sub X1, X1, #%d\n", instruction.value))

		case InstructionOutput:
			// X1 will be cleared after doing the syscall, but we use this as our register for the data pointer.
			// We must save it in X15, and restore it afterwards.
			sb.WriteString("    mov X15, X1\n")

			// X1 = address of message
			// X2 = length
			sb.WriteString("    mov X2, #1\n")
			sb.WriteString("    mov X16, #4\n")
			sb.WriteString("    mov X0, #1\n")
			sb.WriteString("    svc #0x80\n")

			sb.WriteString("    mov X1, X15\n")

		case InstructionInput:
			// X1 will be cleared after doing the syscall, but we use this as our register for the data pointer.
			// We must save it in X15, and restore it afterwards.
			sb.WriteString("    mov X15, X1\n")

			// X0 = fd to read from
			// X1 = address to write to
			// X2 = length
			sb.WriteString("    mov X2, #1\n")
			sb.WriteString("    mov X16, #3\n")
			sb.WriteString("    mov X0, #0\n")
			sb.WriteString("    svc #0x80\n")

			sb.WriteString("    mov X1, X15\n")

		case InstructionJumpIfZero:
			sb.WriteString("    ldrb W0, [X1]\n")
			sb.WriteString(fmt.Sprintf("    cbz W0, _l%d\n", instruction.value))
			sb.WriteString(fmt.Sprintf("_l%d:", index))

		case InstructionJumpIfNonZero:
			sb.WriteString("    ldrb W0, [X1]\n")
			sb.WriteString(fmt.Sprintf("    cbnz W0, _l%d\n", instruction.value))
			sb.WriteString(fmt.Sprintf("_l%d:", index))

		default:
			fmt.Println("ERROR: Unsupported instruction", instruction.instructionType)
			continue
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nexit:\n")
	sb.WriteString("    mov X0, #0\n")
	sb.WriteString("    mov X16, #1\n")
	sb.WriteString("    svc #0x80\n")

	sb.WriteString("\n.data\nCELLS: .zero 30000\n")

	return sb.String()
}

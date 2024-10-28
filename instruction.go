package main

type InstructionType int

const (
	// Increment the data pointer by 1.
	InstructionIncrement = iota

	// Decrement the data pointer by 1.
	InstructionDecrement

	// Increment the byte at the data pointer by 1.
	InstructionAdd

	// Decrement the byte at the data pointer by 1.
	InstructionSubtract

	// Output the byte at the data pointer.
	InstructionOutput

	// Accept one byte of input, storing its value in the byte at the data pointer.
	InstructionInput

	// If the byte at the data pointer is zero, jump to the command after the matching JumpIfNonZero.
	InstructionJumpIfZero

	// If the byte at the data pointer is non-zero, jump to the command after the matching JumpIfZero.
	InstructionJumpIfNonZero
)

var instructionName = map[InstructionType]string{
	InstructionIncrement:     "increment",
	InstructionDecrement:     "decrement",
	InstructionAdd:           "add",
	InstructionSubtract:      "subtract",
	InstructionOutput:        "output",
	InstructionInput:         "input",
	InstructionJumpIfZero:    "jump-if-zero",
	InstructionJumpIfNonZero: "jump-if-non-zero",
}

// Returns a human-readable name for this instruction type.
func (it InstructionType) String() string {
	return instructionName[it]
}

// Represents a single parsed instruction.
type Instruction struct {
	// The type of instruction that this is.
	instructionType InstructionType

	// The associated value with this instruction.
	// - For Increment, Decrement, Add and Subtract instructions, this holds the amount to increment/decrement by.
	// - For JumpX instructions, this holds the index of the instruction to jump to if the condition is true.
	// - For all other instructions, this is unused.
	value int
}

// Continuously reads from the provided lexer until all bytes have been read, returning a list of parsed instructions.
func ParseInstructions(lexer *Lexer) ([]Instruction, error) {
	instructions := []Instruction{}

	// 0 indicates EOF from the lexer, meaning there is nothing left to parse.
	c := lexer.Next()
	for c != 0 {
		var instruction Instruction

		switch c {
		case '+':
			occurrences := 0
			for c == '+' {
				occurrences++
				c = lexer.Next()
			}

			if c != 0 {
				// Next() will consume the character, so if it is not the one we want, we still
				// need to handle it.
				lexer.position -= 1
			}

			instruction = Instruction{instructionType: InstructionIncrement, value: occurrences}
		case '-':
			occurrences := 0
			for c == '-' {
				occurrences++
				c = lexer.Next()
			}

			if c != 0 {
				// Next() will consume the character, so if it is not the one we want, we still
				// need to handle it.
				lexer.position -= 1
			}

			instruction = Instruction{instructionType: InstructionDecrement, value: occurrences}
		case '>':
			occurrences := 0
			for c == '>' {
				occurrences++
				c = lexer.Next()
			}

			if c != 0 {
				// Next() will consume the character, so if it is not the one we want, we still
				// need to handle it.
				lexer.position -= 1
			}

			instruction = Instruction{instructionType: InstructionAdd, value: occurrences}
		case '<':
			occurrences := 0
			for c == '<' {
				occurrences++
				c = lexer.Next()
			}

			if c != 0 {
				// Next() will consume the character, so if it is not the one we want, we still
				// need to handle it.
				lexer.position -= 1
			}

			instruction = Instruction{instructionType: InstructionSubtract, value: occurrences}
		case '.':
			instruction = Instruction{instructionType: InstructionOutput}
		case ',':
			instruction = Instruction{instructionType: InstructionInput}
		case '[':
			// This instruction tells the interpreter to jump to the character after its matching '[' if
			// the byte at the current data pointer is zero.
			// We will need to resolve that location once we have parsed all instructions.
			instruction = Instruction{instructionType: InstructionJumpIfZero}
		case ']':
			instruction = Instruction{instructionType: InstructionJumpIfNonZero}
		}

		instructions = append(instructions, instruction)
		c = lexer.Next()
	}

	// We have now collected all instructions, we must now resolve any jump locations.
	for index, instruction := range instructions {
		if instruction.instructionType == InstructionJumpIfZero {
			position := index
			depth := 1

			// To resolve the jump location, we must traverse through the instructions until we find
			// the matching non-zero instruction.
			// If we come across any if-zero instructions on the way, we must increment our depth, as they will
			// have a matching non-zero instruction before the one we are currently looking for.
			for depth > 0 {
				position++
				peek := instructions[position]

				if peek.instructionType == InstructionJumpIfZero {
					depth++
				} else if peek.instructionType == InstructionJumpIfNonZero {
					depth--
				}
			}

			// We now have the position of the instruction that corresponds with this jump-if-zero.
			instructions[index] = Instruction{instructionType: InstructionJumpIfZero, value: position}
			continue
		}

		if instruction.instructionType == InstructionJumpIfNonZero {
			// We must find the closest if-zero instruction.
			depth := 1
			position := index - 1

			// Once again, to resolve the matching jump-if-zero, we will need to traverse through the instructions.
			// If we come across another jump-if-non-zero, it has a corresponding jump-if-zero before the one
			// we are looking for, so we must keep track of the depth.
			for depth > 0 {
				peek := instructions[position]
				position--

				if peek.instructionType == InstructionJumpIfNonZero {
					depth++
				} else if peek.instructionType == InstructionJumpIfZero {
					depth--
				}
			}

			instructions[index] = Instruction{instructionType: InstructionJumpIfNonZero, value: position + 1}
		}
	}

	return instructions, nil
}

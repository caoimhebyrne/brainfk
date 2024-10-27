package main

type InstructionType int

const (
	Instruction_Inc = iota
	Instruction_Dec
	Instruction_Add
	Instruction_Sub
	Instruction_Output
	Instruction_Input
	Instruction_JumpIfZero
	Instruction_JumpIfNonZero
)

type Instruction struct {
	// The type of instruction that this is.
	instructionType InstructionType

	// The associated value with this instruction.
	// For all instructions except JumpX, this will be the occurrences.
	// For JumpX instructions, this indicates the index of the instruction to jump to
	// if the condition is true.
	value int
}

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

			// Next() will consume the character, so if it is not the one we want, we still
			// need to handle it.
			lexer.position -= 1

			instruction = Instruction{instructionType: Instruction_Inc, value: occurrences}
		case '-':
			occurrences := 0
			for c == '-' {
				occurrences++
				c = lexer.Next()
			}

			// Next() will consume the character, so if it is not the one we want, we still
			// need to handle it.
			lexer.position -= 1

			instruction = Instruction{instructionType: Instruction_Dec, value: occurrences}
		case '>':
			occurrences := 0
			for c == '>' {
				occurrences++
				c = lexer.Next()
			}

			// Next() will consume the character, so if it is not the one we want, we still
			// need to handle it.
			lexer.position -= 1

			instruction = Instruction{instructionType: Instruction_Add, value: occurrences}
		case '<':
			occurrences := 0
			for c == '<' {
				occurrences++
				c = lexer.Next()
			}

			// Next() will consume the character, so if it is not the one we want, we still
			// need to handle it.
			lexer.position -= 1

			instruction = Instruction{instructionType: Instruction_Sub, value: occurrences}
		case '.':
			instruction = Instruction{instructionType: Instruction_Output}
		case ',':
			instruction = Instruction{instructionType: Instruction_Input}
		case '[':
			// This instruction tells the interpreter to jump to the character after its matching '[' if
			// the byte at the current data pointer is zero.
			// We will need to resolve that location once we have parsed all instructions.
			instruction = Instruction{instructionType: Instruction_JumpIfZero}
		case ']':
			instruction = Instruction{instructionType: Instruction_JumpIfNonZero}
		}

		instructions = append(instructions, instruction)
		c = lexer.Next()
	}

	// We have now collected all instructions, we must now resolve any jump locations.
	for index, instruction := range instructions {
		if instruction.instructionType == Instruction_JumpIfZero {
			position := index
			depth := 1

			// To resolve the jump location, we must traverse through the instructions until we find
			// the matching non-zero instruction.
			// If we come across any if-zero instructions on the way, we must increment our depth, as they will
			// have a matching non-zero instruction before the one we are currently looking for.
			for depth > 0 {
				position++
				peek := instructions[position]

				if peek.instructionType == Instruction_JumpIfZero {
					depth++
				} else if peek.instructionType == Instruction_JumpIfNonZero {
					depth--
				}
			}

			// We now have the position of the instruction that corresponds with this jump-if-zero.
			instructions[index] = Instruction{instructionType: Instruction_JumpIfZero, value: position}
			continue
		}

		if instruction.instructionType == Instruction_JumpIfNonZero {
			// We must find the closest if-zero instruction.
			depth := 1
			position := index - 1

			// Once again, to resolve the matching jump-if-zero, we will need to traverse through the instructions.
			// If we come across another jump-if-non-zero, it has a corresponding jump-if-zero before the one
			// we are looking for, so we must keep track of the depth.
			for depth > 0 {
				peek := instructions[position]
				position--

				if peek.instructionType == Instruction_JumpIfNonZero {
					depth++
				} else if peek.instructionType == Instruction_JumpIfZero {
					depth--
				}
			}

			instructions[index] = Instruction{instructionType: Instruction_JumpIfNonZero, value: position}
		}
	}

	return instructions, nil
}

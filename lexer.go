package main

import "os"

type Lexer struct {
	// The unfiltered contents of the input file, may or may not be valid brainfuck.
	contents []byte

	// The current position of the Lexer into the contents.
	position int
}

// Initializes a new lexer with the contents of the file at the provided path.
// Returns the lexer instance if no errors occur when reading the file.
func NewLexer(filePath string) (*Lexer, error) {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Lexer{contents: contents, position: 0}, nil
}

// Returns the next character from the contents.
// This will only return characters supported by the brainfuck language.
// If the end of the stream is reached, 0 will be returned.
func (l *Lexer) Next() byte {
	// When reading bytes from the contents, we want to ignore anything that is not valid brainfuck.
	// To ensure that only valid characters are returned, we will iterate over the array until we reach a valid character.
	for l.position < len(l.contents) {
		character := l.contents[l.position]
		l.position++

		if l.isValidCharacter(character) {
			return character
		}
	}

	return 0
}

// Returns whether the provided character is valid in the brainfuck language.
func (l *Lexer) isValidCharacter(character byte) bool {
	return character == '+' || character == '-' || character == '>' || character == '<' || character == '.' || character == ',' || character == '[' || character == ']'
}

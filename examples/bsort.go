package main

import (
	"fmt"
	"os"
)

func main() {
	input, err := os.ReadFile("./examples/bsort.input")
	if err != nil {
		fmt.Println("ERROR: Failed to read input:", err)
		return
	}

	didSwap := true
	for didSwap {
		didSwap = false
		for i := 0; i < (len(input) - 1); i++ {
			current := input[i]
			next := input[i+1]

			if current > next {
				input[i+1] = current
				input[i] = next
				didSwap = true
			}
		}
	}

	for _, v := range input {
		fmt.Printf("%c", v)
	}
}

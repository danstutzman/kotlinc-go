package main

import (
	"fmt"
	"strings"
)

const SOURCE = `fun main(args: Array<String>) {
  println("Hello, world!")
}`

func main() {
	scanner := NewScanner(strings.NewReader(SOURCE))
	words := []string{}
	for {
		token, matched := scanner.Scan()
		if token == EOF {
			break
		}

		word := token.toString()
		if token == STRING_LITERAL || token == IDENT {
			word = word + "(" + matched + ")"
		}
		words = append(words, word)
	}
	fmt.Println(words)
}

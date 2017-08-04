package main

import (
	"fmt"
	"github.com/danielstutzman/kotlinc-go/parser"
)

//go:generate $GOPATH/bin/pigeon -o parser/kotlin.peg.go parser/kotlin.peg

func main() {
	input := "fun hello(args: Array<String>) { println(\"a\", \"b\") }"
	parsed, err := parser.Parse("", []byte(input))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", parsed)
}

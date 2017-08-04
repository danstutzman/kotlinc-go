package main

import (
	"fmt"
)

//go:generate $GOPATH/bin/pigeon -o kotlin.peg.go kotlin.peg

func main() {
	input := "fun hello(args: Array<String>)"
	parsed, err := Parse("", []byte(input))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", parsed)
}

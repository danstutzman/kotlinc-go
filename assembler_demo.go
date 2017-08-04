package main

import (
	"fmt"
	"github.com/danielstutzman/kotlinc-go/assembler"
	"os"
)

func main() {
	if len(os.Args) < 1+1 {
		fmt.Fprintf(os.Stderr, "Specify .class file to write out\n")
		os.Exit(1)
	}
	outPath := os.Args[1]

	assembler.Demo(outPath)
}

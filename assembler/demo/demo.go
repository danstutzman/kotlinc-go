package main

import (
	"fmt"
	"github.com/danielstutzman/kotlinc-go/assembler"
	"os"
)

func main() {
	classFile := assembler.CreateClassFile("MinimalGo", []string{"a", "b"})

	outPath := "MinimalGo.class"
	out, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	classFile.Write(out)
	fmt.Println(outPath)
}

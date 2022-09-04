package main

import (
	"lsbasi/interpreter"
	"os"
)

func main() {
	text, err := os.ReadFile("./code.goc")
	if err != nil {
		panic(err)
	}

	if len(text) == 0 {
		panic("empty file")
	}

	i := interpreter.New([]rune(string(text)))
	i.Interprete()
	i.PrintValues()
}

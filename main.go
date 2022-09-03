package main

import (
	"bufio"
	"fmt"
	"lsbasi/interpreter"
	"os"
)

func main() {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("calc> ")
		text, _, err := input.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(text) == 0 {
			continue
		}

		i := interpreter.New([]rune(string(text)))
		result := i.Expr()
		fmt.Println(result)
	}
}

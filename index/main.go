package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}
	CreteIndex(os.Args[1])
}

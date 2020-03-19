package main

import (
	"fmt"
	"github.com/polisgo2020/search-senyast4745/index"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}
	index.CreteIndex(os.Args[1])
}

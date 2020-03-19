package main

import (
	"fmt"
	"github.com/polisgo2020/search-senyast4745/index"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("too few arguments")
	}
	index.CreteIndex(os.Args[1])
	index.SearchWordsInIndex(os.Args[1], os.Args[2:])
}

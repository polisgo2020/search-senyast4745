package main

import (
	"fmt"
	utilVocabulary "github.com/senyast4745/index/vocabulary"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}
	utilVocabulary.CreteIndex(os.Args[1])
}

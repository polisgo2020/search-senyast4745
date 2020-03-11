package main

import (
	"fmt"
	"github.com/senyast4745/index/files"
	"github.com/senyast4745/index/util"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}
	util.Check(files.CreteIndex(os.Args[1]), "error %e while creating index")
}

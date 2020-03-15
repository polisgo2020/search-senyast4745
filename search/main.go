package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("too few arguments")
	}
	//util.Check(files.CreteIndex(os.Args[1]), "error %e while creating index")
	if err := ReadCSVFile(os.Args[1], os.Args[2:]); err != nil {
		fmt.Printf("Error %e while parsing CSV file \n", err)
		return
	}
}

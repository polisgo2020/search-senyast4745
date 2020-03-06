package main

import (
	files2 "com.github.senyast4745/index/files"
	"com.github.senyast4745/index/util"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}

	folderLocation := os.Args[1]

	files, err := files2.FilePathWalkDir(folderLocation)
	util.Check(err, "error %e while reading files from directory")
	m := files2.CollectWordData(files)

	js, err := json.Marshal(m)
	util.Check(err, "error %e while making json data")
	util.Check(files2.WriteDataToFile(string(js)), "error %e while saving data to file")

}

package main

import (
	"fmt"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("too few program arguments")
		return
	}
	CreteIndex(os.Args[1])
}

func CreteIndex(folderLocation string) {
	if allFiles, err := files.FilePathWalkDir(folderLocation); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := CollectWordData(allFiles)
		util.Check(index.Index.CollectAndWriteMap(m), "error %e while saving data to file")
	}

}

func CollectWordData(fileNames []string) index.Index {
	m := make(index.Index)
	for fn := range fileNames {

		if words, err := files.ReadFileByWords(fileNames[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, fileNames[fn])
		} else {
			data, err := index.MapAndCleanWords(words, fileNames[fn])
			if err != nil {
				util.Check(err, "error %e")
			}
			for i := range data {
				if m[i] == nil {
					m[i] = []*files.FileStruct{data[i]}
				} else {
					m[i] = append(m[i], data[i])
				}
			}
		}
	}
	return m
}

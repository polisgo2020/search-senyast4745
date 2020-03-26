package main

import (
	"encoding/csv"
	"encoding/json"
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
	creteIndex(os.Args[1])
}

func creteIndex(folderLocation string) {
	if allFiles, err := files.FilePathWalkDir(folderLocation); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := collectWordData(allFiles)
		util.Check(collectAndWriteMap(m), "error %e while saving data to file")
	}

}

func collectWordData(fileNames []string) index.Index {
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

func collectAndWriteMap(ind index.Index) error {
	if err := os.MkdirAll(files.FinalOutputDirectory, 0777); err != nil {
		return err
	}
	recordFile, _ := os.Create(files.FinalDataFile)
	w := csv.NewWriter(recordFile)
	var count int
	for k, v := range ind {
		t, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error %e while creating json from obj %+v \n", err, &v)
		}
		err = w.Write([]string{k, string(t)})
		if err != nil {
			fmt.Printf("error %e while saving record %s,%s \n", err, k, t)
		}
		count++
		if count > 100 {
			w.Flush()
		}
	}
	return nil
}

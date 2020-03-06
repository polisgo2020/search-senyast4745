package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const finalDataFile = "./final.json"

func main() {
	folderLocation := os.Args[1]
	files, err := filePathWalkDir(folderLocation)
	check(err, "error %e while reading files from directory")
	m := make(map[string][]string)
	for fn := range files {
		if words, err := readFileByWords(files[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, files[fn])
		} else {
			for i := range words {
				m[words[i]] = append(m[words[i]], files[fn])
			}
		}
	}

	js, err := json.Marshal(m)
	check(err, "error %e while making json data")
	check(writeDataToFile(string(js)), "error %e while saving data to file")

}

func check(err error, format string) {
	if err != nil {
		fmt.Printf(format, err)
	}
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

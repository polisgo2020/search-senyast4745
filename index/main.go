package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const finalDataFile = "./final.txt"

func main() {
	folderLocation := os.Args[1]
	files, err := filePathWalkDir(folderLocation)
	check(err, "error %e while reading files from directory")
	m := make(map[string][]int)
	for fn := range files {
		if words, err := readFileByWords(files[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, files[fn])
		} else {
			for i := range words {
				m[words[i]] = append(m[words[i]], fn)
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

func readFileByWords(fn string) ([]string, error) {

	file, err := os.Open(fn)
	if err != nil {
		log.Fatalf("error while ")
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func writeDataToFile(str string) error {

	f, err := os.Create(finalDataFile)
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	if _, err := f.WriteString(str); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

package main

import (
	"fmt"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("too few arguments")
		return
	}
	searchWordsInIndex(os.Args[1], os.Args[2:])
}

func searchWordsInIndex(filePath string, words []string) {
	if inputWords, err := util.CleanUserInput(words); err != nil {
		fmt.Printf("Error %e while cleaning user input", err)
	} else {

		data, err := files.ReadCSVFile(filePath)
		if err != nil {
			fmt.Printf("Couldn't open or read the csv file %s with error %e \n", filePath, err)
		}
		for k, v := range getCorrectFiles(data, inputWords) {
			fmt.Printf("Filename: %s, words count: %d, spacing between words in a file: %d \n", k, v.Path, v.Weight)
		}
	}
}

func getCorrectFiles(m index.Index, searchWords []string) map[string]index.Data {
	data := make(index.Index)
	for i := range searchWords {
		data[searchWords[i]] = m[searchWords[i]]
	}
	a := make(index.Index)
	for i := range data {
		dataLen := len(data[i])
		if 0 != dataLen {
			a[i] = data[i]
		}
	}
	return a.Search(searchWords)
}

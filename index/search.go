package index

import (
	"errors"
	"fmt"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/polisgo2020/search-senyast4745/util"
	"github.com/reiver/go-porterstemmer"
	"math"
	"strings"
	"unicode"
)

type data struct {
	file   string
	Weight int
	Path   int
}

func SearchWordsInIndex(filePath string, words []string) {
	if inputWords, err := cleanUserInput(words); err != nil {
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

func getCorrectFiles(m map[string][]*files.FileStruct, searchWords []string) map[string]data {
	data := make(map[string][]*files.FileStruct)
	for i := range searchWords {
		data[searchWords[i]] = m[searchWords[i]]
	}
	a := make(map[string][]*files.FileStruct)
	for i := range data {
		dataLen := len(data[i])
		if 0 != dataLen {
			a[i] = data[i]
		}
	}
	return sortFiles(a, searchWords)
}

func cleanUserInput(words []string) ([]string, error) {
	var data []string
	for _, v := range words {
		word := strings.TrimFunc(v, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if (!EnglishStopWordChecker(word)) && (len(word) > 0) {
			data = append(data, porterstemmer.StemString(word))
		}
	}
	if len(data) == 0 {
		return nil, errors.New("bad search words")
	}
	return data, nil
}

//sorting data by number of occurrences of words and distance between words in the source file
func sortFiles(m map[string][]*files.FileStruct, searchWords []string) map[string]data {
	dataFirst := make(map[int]map[string]data)
	dataSecond := dataFirst
	for i := range searchWords {
		for j := range m[searchWords[i]] {
			for k := range m[searchWords[i]][j].Position {
				minW := math.MaxInt64
				if dataSecond[k] == nil {
					dataSecond[k] = make(map[string]data)
				}
				if _, ok := dataSecond[k][m[searchWords[i]][j].File]; !ok {
					dataSecond[k][m[searchWords[i]][j].File] = data{file: m[searchWords[i]][j].File}
				}
				for t := range dataFirst {
					if dataFirst[t][m[searchWords[i]][j].File].Weight+util.Abs(t-m[searchWords[i]][j].Position[k]) < minW {
						minW = dataFirst[t][m[searchWords[i]][j].File].Weight + util.Abs(t-m[searchWords[i]][j].Position[k])
						dataSecond[t][m[searchWords[i]][j].File] = data{file: m[searchWords[i]][j].File, Weight: minW,
							Path: dataFirst[t][m[searchWords[i]][j].File].Path + 1}
					}
				}
			}
		}
	}
	ans := make(map[string]data)
	for _, v := range dataFirst {
		for k := range v {
			ans[k] = v[k]
		}
	}
	return ans
}

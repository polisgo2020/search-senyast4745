package files

import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"github.com/senyast4745/index/util"
	_vocabulary "github.com/senyast4745/index/vocabulary"
	"strings"
	"unicode"
)

func CreteIndex(folderLocation string) error {
	if allFiles, err := FilePathWalkDir(folderLocation); err != nil {
		return err
	} else {
		m := collectWordData(allFiles)
		return CollectAndWriteMap(m)
	}

}
func mapAndCleanWords(fileData []string, fn string) (map[string]*WordStruct, error) {

	var position int
	data := make(map[string]*WordStruct)
	for i := range fileData {
		word := strings.TrimFunc(fileData[i], func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if (!_vocabulary.EnglishStopWordChecker(word)) && (len(word) > 0) {
			word = porterstemmer.StemString(word)

			if data[word] == nil {
				data[word] = &WordStruct{File: fn, Position: []int{position}}
			} else {
				data[word].Position = append(data[word].Position, position)
			}
			position++
		}
	}
	return data, nil
}

func collectWordData(fileNames []string) map[string][]*WordStruct {
	m := make(map[string][]*WordStruct)
	for fn := range fileNames {

		if words, err := ReadFileByWords(fileNames[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, fileNames[fn])
		} else {
			data, err := mapAndCleanWords(words, fileNames[fn])
			if err != nil {
				util.Check(err, "error %e")
			}
			for i := range data {
				if m[i] == nil {
					m[i] = []*WordStruct{data[i]}
				} else {
					m[i] = append(m[i], data[i])
				}
			}
		}
	}
	return m
}

type WordStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

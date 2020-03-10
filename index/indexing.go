package main

import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"github.com/senyast4745/index/files"
	"github.com/senyast4745/index/util"
	_vocabulary "github.com/senyast4745/index/vocabulary"
	"strings"
	"unicode"
)

func CreteIndex(folderLocation string) {
	if allFiles, err := files.FilePathWalkDir(folderLocation); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := collectWordData(allFiles)
		util.Check(files.CollectAndWriteMap(m), "error %e while saving data to file")
	}

}
func mapAndCleanWords(fileData []string, fn string) (map[string]*_vocabulary.WordStruct, error) {

	var position int
	data := make(map[string]*_vocabulary.WordStruct)
	for i := range fileData {
		word := strings.TrimFunc(fileData[i], func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if (!_vocabulary.EnglishStopWordChecker(word)) && (len(word) > 0) {
			word = porterstemmer.StemString(word)

			if data[word] == nil {
				data[word] = &_vocabulary.WordStruct{File: fn, Position: []int{position}}
			} else {
				data[word].Position = append(data[word].Position, position)
			}
			position++
		}
	}
	return data, nil
}

func collectWordData(fileNames []string) map[string][]*_vocabulary.WordStruct {
	m := make(map[string][]*_vocabulary.WordStruct)
	for fn := range fileNames {

		if words, err := files.ReadFileByWords(fileNames[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, fileNames[fn])
		} else {
			data, err := mapAndCleanWords(words, fileNames[fn])
			if err != nil {
				util.Check(err, "error %e")
			}
			for i := range data {
				if m[i] == nil {
					m[i] = []*_vocabulary.WordStruct{data[i]}
				} else {
					m[i] = append(m[i], data[i])
				}
			}
		}
	}
	return m
}

package index

import (
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/reiver/go-porterstemmer"
	"strings"
	"unicode"
)

// MapAndCleanWords creates an inverted index for a given word slice from a given file
func MapAndCleanWords(fileData []string, fn string) (map[string]*files.FileStruct, error) {
	var position int
	data := make(map[string]*files.FileStruct)
	for i := range fileData {
		word := strings.TrimFunc(fileData[i], func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if !EnglishStopWordChecker(word) {
			word = porterstemmer.StemString(word)
			if len(word) > 0 {
				if data[word] == nil {
					data[word] = &files.FileStruct{File: fn, Position: []int{position}}
				} else {
					data[word].Position = append(data[word].Position, position)
				}
				position++
			}
		}
	}
	return data, nil
}

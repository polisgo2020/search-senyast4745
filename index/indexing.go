package index

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/polisgo2020/search-senyast4745/util"
	"github.com/reiver/go-porterstemmer"
)

// MapAndCleanWords creates an inverted index for a given word slice from a given file

type FileWordMap map[string]*FileStruct

func MapAndCleanWords(reader io.Reader, fn string) (FileWordMap, error) {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanWords)

	var position int
	data := make(FileWordMap)
	for sc.Scan() {
		word := strings.TrimFunc(sc.Text(), func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if !util.EnglishStopWordChecker(word) {
			word = porterstemmer.StemString(word)
			if len(word) > 0 {
				if data[word] == nil {
					data[word] = &FileStruct{File: fn, Position: []int{position}}
				} else {
					data[word].Position = append(data[word].Position, position)
				}
				position++
			}
		}
	}
	return data, nil
}

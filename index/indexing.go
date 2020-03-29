package index

import (
	"bufio"
	"github.com/polisgo2020/search-senyast4745/util"
	"io"
)

type FileWordMap map[string]*FileStruct

// MapAndCleanWords creates an inverted index for a given word slice from a given file
func MapAndCleanWords(reader io.Reader, fn string) (FileWordMap, error) {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanWords)

	var position int
	data := make(FileWordMap)
	for sc.Scan() {
		util.CleanUserInput(sc.Text(), func(input string) {
			if data[input] == nil {
				data[input] = &FileStruct{File: fn, Position: []int{position}}
			} else {
				data[input].Position = append(data[input].Position, position)
			}
			position++
		})
	}
	return data, nil
}

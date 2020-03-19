package index

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/reiver/go-porterstemmer"
	"os"
	"strings"
	"unicode"
)

func MapAndCleanWords(fileData []string, fn string) (map[string]*files.FileStruct, error) {
	var position int
	data := make(map[string]*files.FileStruct)
	for i := range fileData {
		word := strings.TrimFunc(fileData[i], func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if word = stopwords.CleanString(word, "en", true); len(word) > 0 {
			word = porterstemmer.StemString(word)

			if data[word] == nil {
				data[word] = &files.FileStruct{File: fn, Position: []int{position}}
			} else {
				data[word].Position = append(data[word].Position, position)
			}
			position++
		}
	}
	return data, nil
}

func (ind Index) CollectAndWriteMap() error {
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

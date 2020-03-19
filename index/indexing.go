package index

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/polisgo2020/search-senyast4745/util"
	"github.com/reiver/go-porterstemmer"
	"os"
	"strings"
	"unicode"
)

const finalDataFile = "output/final.csv"

const finalOutputDirectory = "output"

func CreteIndex(folderLocation string) {
	if allFiles, err := files.FilePathWalkDir(folderLocation); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := collectWordData(allFiles)
		util.Check(CollectAndWriteMap(m), "error %e while saving data to file")
	}

}
func mapAndCleanWords(fileData []string, fn string) (map[string]*files.FileStruct, error) {

	var position int
	data := make(map[string]*files.FileStruct)
	for i := range fileData {
		word := strings.TrimFunc(fileData[i], func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if (!EnglishStopWordChecker(word)) && (len(word) > 0) {
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

func collectWordData(fileNames []string) map[string][]*files.FileStruct {
	m := make(map[string][]*files.FileStruct)
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
					m[i] = []*files.FileStruct{data[i]}
				} else {
					m[i] = append(m[i], data[i])
				}
			}
		}
	}
	return m
}

func CollectAndWriteMap(m map[string][]*files.FileStruct) error {
	if err := os.MkdirAll(finalOutputDirectory, 0777); err != nil {
		return err
	}
	recordFile, _ := os.Create(finalDataFile)
	w := csv.NewWriter(recordFile)
	var count int
	for k, v := range m {
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

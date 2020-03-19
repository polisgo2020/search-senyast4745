package index

import (
	"fmt"
	"github.com/polisgo2020/search-senyast4745/files"
	"github.com/polisgo2020/search-senyast4745/util"
	"math"
)

type Data struct {
	file   string
	Weight int
	Path   int
}

type FileStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

func SearchWordsInIndex(filePath string, words []string) {
	data, err := files.ReadCSVFile(filePath)
	if err != nil {
		fmt.Printf("Couldn't open or read the csv file %s with error %e \n", filePath, err)
	}
	for k, v := range getCorrectFiles(data, words) {
		fmt.Printf("Filename: %s, words count: %d, spacing between words in a file: %d \n", k, v.Path, v.Weight)
	}
}

func getCorrectFiles(m map[string][]*FileStruct, searchWords []string) map[string]Data {
	data := make(map[string][]*FileStruct)
	for i := range searchWords {
		data[searchWords[i]] = m[searchWords[i]]
	}
	a := make(map[string][]*FileStruct)
	for i := range data {
		dataLen := len(data[i])
		if 0 != dataLen {
			a[i] = data[i]
		}
	}
	return sortFiles(a, searchWords)
}

//sorting data by number of occurrences of words and distance between words in the source file
func sortFiles(m map[string][]*FileStruct, searchWords []string) map[string]Data {
	dataFirst := make(map[int]map[string]Data)
	dataSecond := dataFirst
	for i := range searchWords {
		for j := range m[searchWords[i]] {
			for k := range m[searchWords[i]][j].Position {
				minW := math.MaxInt64
				if dataSecond[k] == nil {
					dataSecond[k] = make(map[string]Data)
				}
				if _, ok := dataSecond[k][m[searchWords[i]][j].File]; !ok {
					dataSecond[k][m[searchWords[i]][j].File] = Data{file: m[searchWords[i]][j].File}
				}
				for t := range dataFirst {
					if dataFirst[t][m[searchWords[i]][j].File].Weight+util.Abs(t-m[searchWords[i]][j].Position[k]) < minW {
						minW = dataFirst[t][m[searchWords[i]][j].File].Weight + util.Abs(t-m[searchWords[i]][j].Position[k])
						dataSecond[t][m[searchWords[i]][j].File] = Data{file: m[searchWords[i]][j].File, Weight: minW,
							Path: dataFirst[t][m[searchWords[i]][j].File].Path + 1}
					}
				}
			}
		}
	}
	ans := make(map[string]Data)
	for _, v := range dataFirst {
		for k := range v {
			ans[k] = v[k]
		}
	}
	return ans
}

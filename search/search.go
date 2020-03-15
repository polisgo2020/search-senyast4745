package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

type data struct {
	file   string
	weight int
	path   int
}

type fileStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

func ReadCSVFile(filePath string, words []string) error {
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Couldn't open the csv file with error %e \n", err)
		return err
	}
	r := csv.NewReader(csvFile)
	data := make(map[string][]*fileStruct)
	var errCount int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error %e while readind csv line \n", err)
			errCount++
			if errCount > 100 {
				return err
			}
			continue
		}
		var tmp []*fileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			fmt.Printf("error %e while parsing json data %s \n", err, record[1])
			continue
		}
		data[record[0]] = tmp
	}
	for k, v := range getFiles(data, words) {
		fmt.Printf("Filename: %s, words count: %d, spacing between words in a file: %d \n", k, v.path, v.weight)
	}
	return nil
}

func getFiles(m map[string][]*fileStruct, searchWords []string) map[string]data {
	data := make(map[string][]*fileStruct)
	for i := range searchWords {
		data[searchWords[i]] = m[searchWords[i]]
	}
	a := make(map[string][]*fileStruct)
	for i := range data {
		dataLen := len(data[i])
		if 0 != dataLen {
			a[i] = data[i]
		}
	}
	return sorFiles(a, searchWords)
}

//sorting data by number of occurrences of words and distance between words in the source file
func sorFiles(m map[string][]*fileStruct, searchWords []string) map[string]data {
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
					if dataFirst[t][m[searchWords[i]][j].File].weight+abs(t-m[searchWords[i]][j].Position[k]) < minW {
						minW = dataFirst[t][m[searchWords[i]][j].File].weight + abs(t-m[searchWords[i]][j].Position[k])
						dataSecond[t][m[searchWords[i]][j].File] = data{file: m[searchWords[i]][j].File, weight: minW,
							path: dataFirst[t][m[searchWords[i]][j].File].path + 1}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/senyast4745/index/files"
	"github.com/senyast4745/index/util"
	"io"
	"math"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("too few arguments")
	}
	util.Check(files.CreteIndex(os.Args[1]), "error %e while creating index")
	csvFile, err := os.Open("output/final.csv")
	if err != nil {
		fmt.Printf("Couldn't open the csv file with error %e", err)
		return
	}
	r := csv.NewReader(csvFile)
	data := make(map[string][]*fileStruct)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error %e while readind csv line", err)
			continue
		}
		var tmp []*fileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			fmt.Printf("error %e while parsing json data %s", err, record[1])
			continue
		}
		data[record[0]] = tmp
	}
	fmt.Printf("%+v \n", getFiles(data, os.Args[2:]))
}

type Data struct {
	file   string
	weight int
	path   int
}

type fileStruct struct {
	File      string `json:"file"`
	Position  []int  `json:"position"`
	frequency int
}

func getFiles(m map[string][]*fileStruct, searchWords []string) map[string]Data {
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
func sorFiles(m map[string][]*fileStruct, searchWords []string) map[string]Data {
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
					if dataFirst[t][m[searchWords[i]][j].File].weight+abs(t-m[searchWords[i]][j].Position[k]) < minW {
						minW = dataFirst[t][m[searchWords[i]][j].File].weight + abs(t-m[searchWords[i]][j].Position[k])
						dataSecond[t][m[searchWords[i]][j].File] = Data{file: m[searchWords[i]][j].File, weight: minW,
							path: dataFirst[t][m[searchWords[i]][j].File].path + 1}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

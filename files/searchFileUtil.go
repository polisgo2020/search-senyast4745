package files

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// FileStruct file structure representing the position of a word in a file
type FileStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

// ReadCSVFile reads a csv file at a given path and converts it to an Index structure
func ReadCSVFile(filePath string) (map[string][]*FileStruct, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	data := make(map[string][]*FileStruct)
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
				return nil, err
			}
			continue
		}
		var tmp []*FileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			fmt.Printf("error %e while parsing json data %s \n", err, record[1])
			continue
		}
		data[record[0]] = tmp
	}
	return data, nil
}

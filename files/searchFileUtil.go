package files

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/polisgo2020/search-senyast4745/index"
	"io"
	"os"
)

func ReadCSVFile(filePath string) (map[string][]*index.FileStruct, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	data := make(map[string][]*index.FileStruct)
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
		var tmp []*index.FileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			fmt.Printf("error %e while parsing json data %s \n", err, record[1])
			continue
		}
		data[record[0]] = tmp
	}
	return data, nil
}

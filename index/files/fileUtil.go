package files

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const finalDataFile = "../output/final.csv"

const finalOutputDirectory = "../output"

func CollectAndWriteMap(m map[string][]*WordStruct) error {
	if err := os.MkdirAll(finalOutputDirectory, 0777); err != nil {
		return err
	}
	recordFile, err := os.Create(finalDataFile)
	if err != nil {
		return err
	}
	w := csv.NewWriter(recordFile)
	defer w.Flush()
	for k, v := range m {
		t, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error %e while creating json from obj %+v \n", err, &v)
		}
		err = w.Write([]string{fmt.Sprintf("%s", k), fmt.Sprintf("%s", string(t))})
		if err != nil {
			fmt.Printf("error %e while saving record %s,%s \n", err, k, t)
		}
	}
	return nil
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ReadFileByWords(fn string) ([]string, error) {

	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

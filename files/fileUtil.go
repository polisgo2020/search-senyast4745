package files

import (
	"bufio"
	"os"
	"path/filepath"
)

const FinalDataFile = "output/final.csv"

const FinalOutputDirectory = "output"

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

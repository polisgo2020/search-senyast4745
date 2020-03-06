package files

import (
	"bufio"
	"com.github.senyast4745/index/vocabulary"
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"os"
	"path/filepath"
	"regexp"
)

const finalDataFile = "./final.json"

func ReadFileByWords(fn string) (map[string]*_vocabulary.WordStruct, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var position int
	data := make(map[string]*_vocabulary.WordStruct)
	for scanner.Scan() {
		word := reg.ReplaceAllString(scanner.Text(), "")
		if (!_vocabulary.EnglishStopWordChecker(word)) && (len(word) > 0) {
			word = porterstemmer.StemString(word)

			if data[word] == nil {
				data[word] = &_vocabulary.WordStruct{File: fn, Position: []int{position}}
			} else {
				data[word].Position = append(data[word].Position, position)
			}
			position++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func WriteDataToFile(str string) error {

	f, err := os.Create(finalDataFile)
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	if _, err := f.WriteString(str); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
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

func CollectWordData(files []string) map[string][]*_vocabulary.WordStruct {
	m := make(map[string][]*_vocabulary.WordStruct)
	for fn := range files {
		if words, err := ReadFileByWords(files[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, files[fn])
		} else {
			for i := range words {
				if m[i] == nil {
					m[i] = []*_vocabulary.WordStruct{words[i]}
				} else {
					m[i] = append(m[i], words[i])
				}
			}
		}
	}
	return m
}

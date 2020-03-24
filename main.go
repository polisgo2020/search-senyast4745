package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Search index"
	app.Usage = "generate index from text files and search over them"

	indexFileFlag := &cli.StringFlag{
		Name:  "index, i",
		Usage: "Index file",
	}

	sourcesFlag := &cli.StringFlag{
		Name:  "sources, s",
		Usage: "Files to index",
	}

	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build search index",
			Flags: []cli.Flag{
				indexFileFlag,
				sourcesFlag,
			},
			Action: build,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search over the index",
			Flags: []cli.Flag{
				indexFileFlag,
			},
			Action: search,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func creteIndex(folderLocation string) {
	if allFiles, err := FilePathWalkDir(folderLocation); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := collectWordData(allFiles)
		util.Check(collectAndWriteMap(m), "error %e while saving data to file")
	}

}

func collectWordData(fileNames []string) *index.Index {
	m := index.NewIndex()
	for fn := range fileNames {

		if words, err := ReadFileByWords(fileNames[fn]); err != nil {
			fmt.Printf("error %e while reading data from file %s", err, fileNames[fn])
		} else {
			data, err := index.MapAndCleanWords(words, fileNames[fn])
			if err != nil {
				util.Check(err, "error %e")
			}
			for i := range data {
				if m.Data[i] == nil {
					m.Data[i] = []*index.FileStruct{data[i]}
				} else {
					m.Data[i] = append(m.Data[i], data[i])
				}
			}
		}
	}
	return m
}

func collectAndWriteMap(ind index.Index) error {
	if err := os.MkdirAll(FinalOutputDirectory, 0777); err != nil {
		return err
	}
	recordFile, _ := os.Create(FinalDataFile)
	w := csv.NewWriter(recordFile)
	var count int
	for k, v := range ind.Data {
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

func searchWordsInIndex(filePath string, words []string) {
	if inputWords, err := util.CleanUserInput(words); err != nil {
		fmt.Printf("Error %e while cleaning user input", err)
	} else {

		data, err := ReadCSVFile(filePath)
		if err != nil {
			fmt.Printf("Couldn't open or read the csv file %s with error %e \n", filePath, err)
		}
		for k, v := range getCorrectFiles(data, inputWords) {
			fmt.Printf("Filename: %s, words count: %d, spacing between words in a file: %d \n", k, v.Path, v.Weight)
		}
	}
}

func getCorrectFiles(m *index.Index, searchWords []string) map[string]*index.Data {
	data := index.NewIndex()
	for i := range searchWords {
		tmp := m.Data[searchWords[i]]
		if len(tmp) != 0 {
			data.Data[searchWords[i]] = tmp
		}
	}
	return data.Search(searchWords)
}

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

const FinalDataFile = "output/final.csv"

const FinalOutputDirectory = "output"

// FilePathWalkDir bypasses the given director and returns a list of all files in this folder
// and returns an error if it is not possible to access the folde
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

// ReadFileByWords reads the given file by words and returns an array of the layer or an error if it is impossible to open or read the file
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

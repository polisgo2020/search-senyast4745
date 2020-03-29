package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/polisgo2020/search-senyast4745/index"
	"github.com/polisgo2020/search-senyast4745/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Search index"
	app.Usage = "generate index from text files and search over them"

	indexFileFlag := &cli.StringFlag{
		Aliases: []string{"i"},
		Name:    "index",
		Usage:   "Index file",
	}

	sourcesFlag := &cli.StringFlag{
		Aliases: []string{"s"},
		Name:    "sources, s",
		Usage:   "Files to index",
	}

	searchFlag := &cli.StringFlag{
		Aliases: []string{"sw"},
		Name:    "search-word, sw",
		Usage:   "Search words separated by comma",
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
				searchFlag,
			},
			Action: search,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func build(c *cli.Context) error {
	if err := checkFlags(c, "index", "sources"); err != nil {
		fmt.Printf("Error %e while checking context", err)
		return nil
	}
	if allFiles, err := filePathWalkDir(c.String("sources")); err != nil {
		util.Check(err, "error %e while reading files from directory")
	} else {
		m := collectWordData(allFiles)
		util.Check(collectAndWriteMap(m, c.String("index")), "error %e while saving data to file")
	}
	return nil
}

func collectWordData(fileNames []string) *index.Index {
	m := index.NewIndex()
	var wg sync.WaitGroup
	for i := range fileNames {
		wg.Add(1)
		go readFileByWords(&wg, m.DataChannel, fileNames[i])
	}

	go func(wg *sync.WaitGroup, readChan chan index.FileWordMap) {
		wg.Wait()
		close(readChan)
	}(&wg, m.DataChannel)
	for data := range m.DataChannel {
		for j := range data {
			if m.Data[j] == nil {
				m.Data[j] = []*index.FileStruct{data[j]}
			} else {
				m.Data[j] = append(m.Data[j], data[j])
			}
		}
	}

	return m
}

func collectAndWriteMap(ind *index.Index, indexFile string) error {
	recordFile, _ := os.Create(indexFile)
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

func search(c *cli.Context) error {
	if err := checkFlags(c, "index", "search-word"); err != nil {
		fmt.Printf("Error %e while checking context", err)
		return nil
	}

	inputWords := make([]string, 0)
	for _, word := range strings.Split(c.String("search-word"), ",") {
		util.CleanUserInput(word, func(input string) {
			inputWords = append(inputWords, input)
		})
	}

	if len(inputWords) == 0 {
		fmt.Printf("Incorrect search words")
		return nil
	}

	data, err := readCSVFile(c.String("index"))
	if err != nil {
		fmt.Printf("Couldn't open or read the csv file %s with error %e \n", c.String("path"), err)
	}
	for k, v := range getCorrectFiles(data, inputWords) {
		fmt.Printf("Filename: %s, words count: %d, spacing between words in a file: %d \n", k, v.Path, v.Weight)
	}

	return nil
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

func readCSVFile(filePath string) (*index.Index, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	data := index.NewIndex()
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
		data.Data[record[0]] = tmp
	}
	return data, nil
}

// FilePathWalkDir bypasses the given director and returns a list of all files in this folder
// and returns an error if it is not possible to access the folder
func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ReadFileByWords reads the given file by words and returns an array of the layer
//or an error if it is impossible to open or read the file
func readFileByWords(wg *sync.WaitGroup, outputCh chan<- index.FileWordMap, fn string) {
	defer wg.Done()
	file, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Error %e while openig file %s", err, fn)
		return
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	if wordMap, err := index.MapAndCleanWords(file, fn); err != nil {
		fmt.Printf("error %e while indexing file %s", err, fn)
	} else {
		outputCh <- wordMap
	}
	return
}

func checkFlags(c *cli.Context, str ...string) error {
	for _, flag := range str {
		if c.String(flag) == "" {
			return errors.New(fmt.Sprintf("empty flag %s", flag))
		}
	}
	return nil
}

package index

import (
	"bufio"
	"io"
	"sync"

	"github.com/polisgo2020/search-senyast4745/util"
)

type FileStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

type FileWordMap map[string]*FileStruct

type Index struct {
	Data        map[string][]*FileStruct
	dataChannel chan FileWordMap
}

func NewIndex() *Index {
	return &Index{Data: make(map[string][]*FileStruct)}
}

func (ind *Index) add(word string, data []*FileStruct) {
	ind.Data[word] = data
}

func (ind *Index) OpenApplyAndListenChannel(consumer func(wg *sync.WaitGroup)) {
	ind.dataChannel = make(chan FileWordMap, 1000)
	var wg sync.WaitGroup
	consumer(&wg)

	go func(wg *sync.WaitGroup, readChan chan FileWordMap) {
		wg.Wait()
		close(readChan)
	}(&wg, ind.dataChannel)

	for data := range ind.dataChannel {
		for j := range data {
			if ind.Data[j] == nil {
				ind.Data[j] = []*FileStruct{data[j]}
			} else {
				ind.Data[j] = append(ind.Data[j], data[j])
			}
		}
	}
}

// MapAndCleanWords creates an inverted index for a given word slice from a given file
func (ind *Index) MapAndCleanWords(reader io.Reader, fn string) {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanWords)

	var position int
	data := make(FileWordMap)
	for sc.Scan() {
		util.CleanUserInput(sc.Text(), func(input string) {
			if data[input] == nil {
				data[input] = &FileStruct{File: fn, Position: []int{position}}
			} else {
				data[input].Position = append(data[input].Position, position)
			}
			position++
		})
	}
	ind.dataChannel <- data
}

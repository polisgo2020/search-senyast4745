package index

import (
	"math"
	"sync"

	"github.com/polisgo2020/search-senyast4745/util"
)

type Data struct {
	file   string
	Weight int
	Path   int
}

type FileStruct struct {
	File     string `json:"file"`
	Position []int  `json:"position"`
}

type Index struct {
	Data        map[string][]*FileStruct
	dataChannel chan FileWordMap
}

func NewIndex() *Index {
	return &Index{Data: make(map[string][]*FileStruct)}
}

func (ind *Index) OpenChannel() {
	ind.dataChannel = make(chan FileWordMap, 1000)
}

func (ind *Index) Listen(wg *sync.WaitGroup) {
	go func(wg *sync.WaitGroup, readChan chan FileWordMap) {
		wg.Wait()
		close(readChan)
	}(wg, ind.dataChannel)

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

// Search sorting Index data by number of occurrences of words and distance between words in the source file
func (ind *Index) Search(searchWords []string) map[string]*Data {
	dataFirst := make(map[int]map[string]*Data)
	dataSecond := dataFirst
	for i := range searchWords {
		for j := range ind.Data[searchWords[i]] {
			for k := range ind.Data[searchWords[i]][j].Position {
				minW := math.MaxInt64
				if dataSecond[k] == nil {
					dataSecond[k] = make(map[string]*Data)
				}
				if _, ok := dataSecond[k][ind.Data[searchWords[i]][j].File]; !ok {
					dataSecond[k][ind.Data[searchWords[i]][j].File] = &Data{file: ind.Data[searchWords[i]][j].File}
				}
				for t := range dataFirst {
					if dataFirst[t][ind.Data[searchWords[i]][j].File].Weight+util.Abs(t-ind.Data[searchWords[i]][j].Position[k]) < minW {
						minW = dataFirst[t][ind.Data[searchWords[i]][j].File].Weight + util.Abs(t-ind.Data[searchWords[i]][j].Position[k])
						dataSecond[t][ind.Data[searchWords[i]][j].File] = &Data{file: ind.Data[searchWords[i]][j].File, Weight: minW,
							Path: dataFirst[t][ind.Data[searchWords[i]][j].File].Path + 1}
					}
				}
			}
		}
	}
	ans := make(map[string]*Data)
	for _, v := range dataFirst {
		for k := range v {
			ans[k] = v[k]
		}
	}
	return ans
}

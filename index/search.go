package index

import (
	"github.com/polisgo2020/search-senyast4745/util"
	"math"
	"sort"
)

type Data struct {
	Weight int
	Path   int
}

type dynamicData struct {
	Path  int
	DPVar []*dynamicVar
}

type dynamicVar struct {
	Position int
	Weight   int
}

func makeDynamicVar(pos []int) []*dynamicVar {
	var t []*dynamicVar
	for i := range pos {
		t = append(t, &dynamicVar{Position: pos[i]})
	}
	return t
}

// Search sorting Index data by number of occurrences of words and distance between words in the source file
// use dynamic programming as search algorithm
func (ind *Index) Search(searchWords []string) map[string]*Data {

	data := make(map[string]*dynamicData)
	for _, word := range searchWords {
		for _, fileStr := range ind.Data[word] {
			if data[fileStr.File] == nil {
				data[fileStr.File] = &dynamicData{DPVar: makeDynamicVar(fileStr.Position)}
			} else {
				data[fileStr.File].DPVar = dynamicMinPosition(data[fileStr.File].DPVar, fileStr.Position)
			}
			data[fileStr.File].Path++
		}
	}
	res := make(map[string]*Data)
	for s := range data {
		res[s] = transform(data[s])
	}
	return res
}

func dynamicMinPosition(dp []*dynamicVar, pos []int) []*dynamicVar {
	for v := range dp {
		dp[v].Weight += findMinDiffPos(pos, dp[v].Position)
	}
	return dp
}

func findMinDiffPos(pos []int, key int) int {
	i := sort.SearchInts(pos, key)
	var diff int
	if i == 0 {
		diff = util.Abs(key - pos[i])
	} else {
		if i == len(pos) {
			diff = util.Abs(key - pos[i-1])
		} else {
			t := util.Abs(key - pos[i])
			k := util.Abs(key - pos[i-1])
			if t < k {
				diff = t
			} else {
				diff = k
			}
		}
	}
	return diff
}

func transform(dd *dynamicData) *Data {
	data := &Data{Path: dd.Path}
	min := math.MaxInt32
	for i := range dd.DPVar {
		if dd.DPVar[i].Weight < min {
			min = dd.DPVar[i].Weight
		}
	}
	data.Weight = min
	return data
}

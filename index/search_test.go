package index

import (
	"github.com/stretchr/testify/require"
	//"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"
)

type searchTestSuite struct {
	suite.Suite
	index *Index
}

func (i *searchTestSuite) SetupTest() {
	i.index = NewIndex()
	FillDefaultIndex(i.index)
}

func FillDefaultIndex(i *Index) {
	i.Data["hello"] = []*FileStruct{{
		File:     "file1",
		Position: []int{0, 5},
	}, {
		File:     "file2",
		Position: []int{6},
	}}
	i.Data["world"] = []*FileStruct{{
		File:     "file1",
		Position: []int{3},
	}, {
		File:     "file2",
		Position: []int{0},
	}, {
		File:     "file3",
		Position: []int{3},
	},
	}
	i.Data["golang"] = []*FileStruct{{
		File:     "file2",
		Position: []int{4, 11},
	}, {
		File:     "file3",
		Position: []int{6},
	}}
}

func (i *searchTestSuite) TestIndex_SimpleSearch() {

	require.Equal(i.T(), 2, len(i.index.Search([]string{"hello"})))
}

func (i *searchTestSuite) TestIndex_Search() {

	expected := map[string]*Data{
		"file1": {
			Weight: 2,
			Path:   2,
		},
		"file2": {
			Weight: 6,
			Path:   2,
		},
		"file3": {
			Weight: 0,
			Path:   1,
		},
	}

	require.Equal(i.T(), expected, i.index.Search([]string{"hello", "world"}))
}

func (i *searchTestSuite) TestIndex_Search2() {

	expected := map[string]*Data{
		"file2": {
			Weight: 0,
			Path:   1,
		},
		"file3": {
			Weight: 0,
			Path:   1,
		},
	}

	require.Equal(i.T(), expected, i.index.Search([]string{"golang"}))
}

func TestSearchSuitStart(t *testing.T) {
	suite.Run(t, new(searchTestSuite))
}

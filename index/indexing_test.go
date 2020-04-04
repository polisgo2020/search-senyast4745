package index

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNewIndex(t *testing.T) {
	ind := NewIndex()
	require.NotNil(t, ind, "check index is not null")
	require.Nil(t, ind.dataChannel)
	require.NotNil(t, ind.Data)
	require.True(t, isClosed(ind.dataChannel), "channel must be closed")
}

type indexTestSuite struct {
	suite.Suite
	called        bool
	index         *Index
	fileWorldMaps []fileWordMap
	fileNames     []string
	input         string
}

func (i *indexTestSuite) SetupTest() {
	i.index = NewIndex()
	i.index.dataChannel = make(chan fileWordMap, 10)

	i.fileNames = []string{"file1", "file2"}

	for _, name := range i.fileNames {
		i.fileWorldMaps = append(i.fileWorldMaps, setupDataToFileMap(name))
	}

	i.called = false

	i.input = "hello world hello Golang Golang"

}

func setupDataToFileMap(filename string) fileWordMap {

	fileWordMap := make(fileWordMap)
	fileWordMap["hello"] = &FileStruct{
		File: filename, Position: []int{0, 2},
	}
	fileWordMap["world"] = &FileStruct{
		File: filename, Position: []int{1},
	}
	fileWordMap["golang"] = &FileStruct{
		File: filename, Position: []int{3, 4},
	}
	return fileWordMap
}

func (i *indexTestSuite) TearDownTest() {
	i.fileWorldMaps = make([]fileWordMap, 0, 3)
	i.fileNames = make([]string, 0, 3)
	if !isClosed(i.index.dataChannel) {
		close(i.index.dataChannel)
	}
}

func TestIndexingSuitStart(t *testing.T) {
	suite.Run(t, new(indexTestSuite))
}

func (i *indexTestSuite) TestIndex_MapAndCleanWords_WithoutStopWords() {
	i.index.MapAndCleanWords(bytes.NewBufferString(i.input), i.fileNames[0])
	actual := <-i.index.dataChannel
	require.Equal(i.T(), i.fileWorldMaps[0], actual)
}

func (i *indexTestSuite) TestIndex_MapAndCleanWords_WithStopWords() {
	i.input += " you are"
	i.index.MapAndCleanWords(bytes.NewBufferString(i.input), i.fileNames[1])
	actual := <-i.index.dataChannel
	require.Equal(i.T(), i.fileWorldMaps[1], actual)
}

func (i *indexTestSuite) TestIndex_MapAndCleanWords_WithNewData() {
	i.input += " world"
	i.fileWorldMaps[0]["world"].Position = append(i.fileWorldMaps[0]["world"].Position, 5)
	i.index.MapAndCleanWords(bytes.NewBufferString(i.input), i.fileNames[0])
	actual := <-i.index.dataChannel
	require.Equal(i.T(), i.fileWorldMaps[0], actual)
}

func (i *indexTestSuite) TestIndex_MapAndCleanWords_WithNewData2() {
	i.input += " test"
	i.fileWorldMaps[0]["test"] = &FileStruct{
		File: i.fileNames[0], Position: []int{5},
	}
	i.index.MapAndCleanWords(bytes.NewBufferString(i.input), i.fileNames[0])
	actual := <-i.index.dataChannel
	require.Equal(i.T(), i.fileWorldMaps[0], actual)
}

func (i *indexTestSuite) TestIndex_SimpleOpenApplyAndListenChannel() {
	i.index.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {
		i.called = true
	})
	require.True(i.T(), i.called)
	require.Empty(i.T(), i.index.Data, "data must not be written")
	require.True(i.T(), isClosed(i.index.dataChannel))
}

func (i *indexTestSuite) TestIndex_OpenApplyAndListenChannel1() {
	i.index.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {
		i.called = true
		i.index.dataChannel <- i.fileWorldMaps[0]
	})
	require.True(i.T(), i.called)
	require.NotEmpty(i.T(), i.index.Data, "data must not be written")
	require.Equal(i.T(), len(i.fileWorldMaps[0]), len(i.index.Data))
	require.Equal(i.T(), i.fileWorldMaps[0]["hello"], i.index.Data["hello"][0])

	_, ok := i.index.Data["test"]
	require.False(i.T(), ok)
}

func (i *indexTestSuite) TestIndex_OpenApplyAndListenChannel2() {
	i.index.OpenApplyAndListenChannel(func(wg *sync.WaitGroup) {
		i.called = true
		i.index.dataChannel <- i.fileWorldMaps[0]
		i.index.dataChannel <- i.fileWorldMaps[1]
	})
	require.True(i.T(), i.called)
	require.NotEmpty(i.T(), i.index.Data, "data must not be written")
	require.Equal(i.T(), len(i.fileWorldMaps[0]), len(i.index.Data))
	require.Equal(i.T(), 2, len(i.index.Data["hello"]))
}

func isClosed(ch <-chan fileWordMap) bool {
	if ch == nil {
		return true
	}
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

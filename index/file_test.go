package index

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"reflect"
	"sync"
	"testing"
)

type fileTestSuite struct {
	suite.Suite
	index           *Index
	defaultStrIndex string
	wbuffer         *bytes.Buffer
	rbuffer         *bytes.Buffer
	decoder         *CsvDecoder
	encoder         *CsvEncoder
}

func TestFileSuitStart(t *testing.T) {
	suite.Run(t, new(fileTestSuite))
}

func (f *fileTestSuite) SetupTest() {

	f.index = NewIndex()
	FillDefaultIndex(f.index)

	f.defaultStrIndex = `hello,"[{""file"":""file1"",""position"":[0,5]},{""file"":""file2"",""position"":[6]}]"
world,"[{""file"":""file1"",""position"":[3]},{""file"":""file2"",""position"":[0]},{""file"":""file3"",""position"":[3]}]"
golang,"[{""file"":""file2"",""position"":[4,11]},{""file"":""file3"",""position"":[6]}]"
`
	f.wbuffer = &bytes.Buffer{}
	f.rbuffer = bytes.NewBufferString(f.defaultStrIndex)

	f.decoder = NewCsvDecoder(f.rbuffer)
	f.encoder = NewCsvEncoder(f.wbuffer)

}

func (f *fileTestSuite) SimpleIndex() {
	f.index = NewIndex()
	f.index.Data["hello"] = []*FileStruct{{
		File:     "file1",
		Position: []int{0, 5},
	},
	}

	f.defaultStrIndex = `hello,"[{""file"":""file1"",""position"":[0,5]}]"` + "\n"
}

func (f *fileTestSuite) TestCsvDecoder_DecodeJson() {
	ind := NewIndex()
	assert.Nil(f.T(), f.decoder.DecodeJson(ind))
	assert.True(f.T(), reflect.DeepEqual(ind, f.index))
}

func (f *fileTestSuite) TestCsvEncoder_EncodeJson() {
	f.SimpleIndex()

	err := f.encoder.EncodeJson(f.index)
	assert.Equal(f.T(), f.wbuffer, bytes.NewBufferString(f.defaultStrIndex))
	assert.Nil(f.T(), err)
}

func (f *fileTestSuite) TestIndex_FromFile() {
	err := f.index.FromFile(f.decoder)
	ind := NewIndex()
	FillDefaultIndex(ind)
	assert.True(f.T(), reflect.DeepEqual(ind, f.index))
	assert.Nil(f.T(), err)
}

func (f *fileTestSuite) TestIndex_ToFile() {
	f.SimpleIndex()

	err := f.index.ToFile(f.encoder)
	assert.Equal(f.T(), f.wbuffer, bytes.NewBufferString(f.defaultStrIndex))
	assert.Nil(f.T(), err)
}

func TestNewCsvDecoder(t *testing.T) {
	tests := []struct {
		name       string
		wantWriter string
		want       *CsvEncoder
	}{
		{name: "new default csv decoder",
			wantWriter: "",
			want: &CsvEncoder{
				m:      &sync.RWMutex{},
				writer: &bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			got := NewCsvEncoder(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("NewCsvDecoder() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCsvDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCsvEncoder(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want *CsvDecoder
	}{
		{name: "new default csv encoder",
			args: args{
				reader: &bytes.Buffer{},
			},
			want: &CsvDecoder{
				m:      &sync.RWMutex{},
				reader: &bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCsvDecoder(tt.args.reader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCsvEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

package index

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

type TestFileData struct {
	data string
}

func (m *TestFileData) ToString() string {
	return m.data
}

func (m *TestFileData) FromString(str string) {
	m.data = str
}

var testConstructor = func() FileData {
	return &TestFileData{}
}

func (f *fileTestSuite) TestCsvDecoder_DecodeJson() {
	dataChannel := make(chan []FileData, 10)

	actual := make(map[int][]FileData)
	expected := make(map[int][]FileData)

	expected[0] = []FileData{newSimpleFileData("hello"),
		newSimpleFileData(`"[{""file"":""file1"",""position"":[0,5]},{""file"":""file2"",""position"":[6]}]"`)}
	expected[1] = []FileData{newSimpleFileData("world"),
		newSimpleFileData(`"[{""file"":""file1"",""position"":[3]},{""file"":""file2"",""position"":[0]},{""file"":""file3"",""position"":[3]}]"`)}
	expected[2] = []FileData{newSimpleFileData("golang"),
		newSimpleFileData(`"[{""file"":""file2"",""position"":[4,11]},{""file"":""file3"",""position"":[6]}]"`)}
	asyncR := func(dataCh <-chan []FileData) {
		var pos int
		for d := range dataCh {
			actual[pos] = d
		}
	}
	go asyncR(dataChannel)
	assert.NoError(f.T(), f.decoder.Decode(dataChannel, testConstructor))
	assert.Panics(f.T(), func() {
		close(dataChannel)
	}, "channel must be closed after decoding")
}

func (f *fileTestSuite) TestCsvEncoder_EncodeJson() {

	var asyncW = func(dataCh chan<- []FileData, data []FileData) {
		for i := 0; i < 9; i++ {
			var t []FileData
			t = append(t, data...)
			dataCh <- t
		}
		assert.Equal(f.T(), 0, f.wbuffer.Len())
		for i := 0; i < 11; i++ {
			var t []FileData
			t = append(t, data...)
			dataCh <- t
		}

		assert.NotEqual(f.T(), 0, f.wbuffer.Len())

		close(dataCh)
	}

	dataChannel := make(chan []FileData, 10)

	a := testConstructor()
	b := testConstructor()
	a.FromString("a")
	b.FromString("b")
	go asyncW(dataChannel, []FileData{a, b})

	assert.NoError(f.T(), f.encoder.Encode(dataChannel))
	assert.Equal(f.T(), bytes.NewBufferString(strings.Repeat("a,b\n", 20)), f.wbuffer)

	f.wbuffer.Reset()

	dataChannel = make(chan []FileData, 10)

	a.FromString("hello")
	b.FromString("golang")
	c := &TestFileData{data: "world"}

	go asyncW(dataChannel, []FileData{a, b, c})

	assert.NoError(f.T(), f.encoder.Encode(dataChannel))
	assert.Equal(f.T(), bytes.NewBufferString(strings.Repeat("hello,golang,world\n", 20)), f.wbuffer)

	f.wbuffer.Reset()

	dataChannel = make(chan []FileData, 10)

	go func(dataCh chan<- []FileData, data []FileData) {
		for i := 0; i < 9; i++ {
			var t []FileData
			t = append(t, data...)
			dataCh <- t
		}
		assert.Equal(f.T(), 0, f.wbuffer.Len())
		close(dataCh)
	}(dataChannel, []FileData{a, b})

	assert.NoError(f.T(), f.encoder.Encode(dataChannel))
	assert.Equal(f.T(), bytes.NewBufferString(strings.Repeat("hello,golang\n", 9)), f.wbuffer)

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

func TestNewSimpleFileData(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want FileData
	}{
		{
			name: "test constructor",
			args: args{
				str: "test",
			},
			want: &simpleFileData{
				data: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSimpleFileData(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleFileData_ToString(t *testing.T) {
	type fields struct {
		data string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: " test to string",
			fields: fields{
				data: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &simpleFileData{
				data: tt.fields.data,
			}
			if got := f.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleFileData_FromString(t *testing.T) {
	type fields struct {
		data string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test from string",
			fields: fields{
				data: "test",
			},
			args: args{
				str: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &simpleFileData{
				data: tt.fields.data,
			}
			f.FromString(tt.args.str)
			assert.Equal(t, tt.args.str, f.data)
		})
	}
}

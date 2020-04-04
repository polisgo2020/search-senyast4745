package index

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

// FromFile with the help of a given decoder reads and decodes the index file and translates it into an index structure
func (ind *Index) FromFile(decoder Decoder) error {

	dataChannel := make(chan []FileData, 10)

	go func(dataCh <-chan []FileData) {
		for data := range dataCh {
			var tmp []*FileStruct
			if err := json.Unmarshal([]byte(data[1].ToString()), &tmp); err != nil {
				log.Println("error", err,
					"msg", fmt.Sprintf("can not parse json data %s \n", data[1].ToString()), "data")
				continue
			}
			ind.add(data[0].ToString(), tmp)
		}
	}(dataChannel)

	return decoder.Decode(dataChannel, func() FileData {
		return &simpleFileData{}
	})
}

// ToFile using the specified encoder saves data to the specified writer
func (ind *Index) ToFile(encoder Encoder) error {

	dataChannel := make(chan []FileData, 10)

	go func(dataCh chan<- []FileData) {
		for i := range ind.Data {
			rawData, err := json.Marshal(ind.Data[i])
			if err != nil {
				log.Printf("Error %q while marshalling data %+v", err, ind.Data[i])
				continue
			}
			dataCh <- []FileData{newSimpleFileData(i), newSimpleFileData(string(rawData))}
		}
		close(dataCh)

	}(dataChannel)

	return encoder.Encode(dataChannel)
}

// CsvDecoder structure for reading and decoding csv file index
type CsvDecoder struct {
	m      *sync.RWMutex
	reader io.Reader
}

// CsvEncoder structure for writing an index to a csv file
type CsvEncoder struct {
	m      *sync.RWMutex
	writer io.Writer
}

// Encode saves data from a channel to a csv file
func (c *CsvEncoder) Encode(dataChannel <-chan []FileData) error {

	w := csv.NewWriter(c.writer)
	defer func() {
		c.m.Lock()
		w.Flush()
		c.m.Unlock()
	}()

	var count int
	for data := range dataChannel {

		var csvData []string

		for i := range data {
			csvData = append(csvData, data[i].ToString())
		}

		err := w.Write(csvData)
		if err != nil {
			log.Printf("can not save record %v \n", csvData)
			return err
		}
		count++
		if count > 10 {
			c.m.Lock()
			w.Flush()
			c.m.Unlock()

			log.Println("msg", "flush writer", "writer", w)
			count = 0
		}
	}

	return nil
}

// Decode reads line-by-line data from a csv file and writes it to the channel
func (c *CsvDecoder) Decode(dataChannel chan<- []FileData, constructor func() FileData) error {
	r := csv.NewReader(c.reader)
	var errCount int
	defer close(dataChannel)

	for {
		c.m.RLock()
		record, err := r.Read()
		c.m.RUnlock()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error", err,
				"msg", "can not read csv line")
			errCount++
			if errCount > 100 {
				return err
			}
			continue
		}
		log.Println("msg", "reading data from csv", "data", record)

		var rawData []FileData

		for i := range record {
			d := constructor()
			d.FromString(record[i])
			rawData = append(rawData, d)
		}
		dataChannel <- rawData
	}
	return nil
}

// NewCsvDecoder default constructor to CsvDecoder with reader
func NewCsvDecoder(reader io.Reader) *CsvDecoder {
	return &CsvDecoder{m: &sync.RWMutex{}, reader: reader}
}

// NewCsvEncoder default constructor to CsvEncoder with writer
func NewCsvEncoder(writer io.Writer) *CsvEncoder {
	return &CsvEncoder{m: &sync.RWMutex{}, writer: writer}
}

type Decoder interface {
	Decode(chan<- []FileData, func() FileData) error
}

type Encoder interface {
	Encode(<-chan []FileData) error
}

// FileData interface for data transfer and processing
type FileData interface {
	FromString(str string)
	ToString() string
}

type simpleFileData struct {
	data string
}

func (f *simpleFileData) ToString() string {
	return f.data
}

func (f *simpleFileData) FromString(str string) {
	f.data = str
}

func newSimpleFileData(str string) FileData {
	return &simpleFileData{data: str}
}

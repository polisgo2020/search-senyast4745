package index

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

func (ind *Index) FromFile(decoder JsonDecoder) error {
	return decoder.DecodeJson(ind)
}

func (ind *Index) ToFile(encoder JsonEncoder) error {
	return encoder.EncodeJson(ind)
}

type CsvDecoder struct {
	m      *sync.RWMutex
	reader io.Reader
}

type CsvEncoder struct {
	m      *sync.RWMutex
	writer io.Writer
}

func (c *CsvEncoder) EncodeJson(ind *Index) error {

	w := csv.NewWriter(c.writer)
	var count int
	for k, v := range ind.Data {
		t, err := json.Marshal(v)
		if err != nil {
			log.Printf("can not create json from obj %+v \n", &v)
			return err
		}
		err = w.Write([]string{k, string(t)})
		if err != nil {
			log.Printf("can not save record %s,%s \n", k, t)
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
	w.Flush()

	return nil
}

func (c *CsvDecoder) DecodeJson(data *Index) error {
	r := csv.NewReader(c.reader)
	var errCount int

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
		var tmp []*FileStruct
		if json.Unmarshal([]byte(record[1]), &tmp) != nil {
			log.Println("error", err,
				"msg", fmt.Sprintf("can not parse json data %s \n", record[1]), "data", record[1])
			continue
		}
		data.add(record[0], tmp)
	}
	return nil
}

func NewCsvDecoder(reader io.Reader) *CsvDecoder {
	return &CsvDecoder{m: &sync.RWMutex{}, reader: reader}
}

func NewCsvEncoder(writer io.Writer) *CsvEncoder {
	return &CsvEncoder{m: &sync.RWMutex{}, writer: writer}
}

type JsonDecoder interface {
	DecodeJson(i *Index) error
}

type JsonEncoder interface {
	EncodeJson(i *Index) error
}

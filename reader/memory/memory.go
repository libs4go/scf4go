package memory

import (
	"encoding/json"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/libs4go/errors"
	"github.com/libs4go/scf4go"
)

type memoryReader struct {
	options []Option
	blocks  []*scf4go.ReadBlock
	sj      *simplejson.Json
}

// ReaderWriter .
type ReaderWriter interface {
	scf4go.Reader
	Write(value interface{}, path ...string)
}

// Option .
type Option func(reader *memoryReader) error

// New .
func New(options ...Option) ReaderWriter {

	reader := &memoryReader{
		options: options,
		sj:      simplejson.New(),
	}

	return reader
}

func (reader *memoryReader) Read() ([]*scf4go.ReadBlock, error) {

	data, err := reader.sj.Encode()

	if err != nil {
		return nil, errors.Wrap(err, "marshal json error")
	}

	reader.blocks = append(reader.blocks, &scf4go.ReadBlock{
		Data:      []byte(data),
		Codec:     "json",
		Timestamp: time.Now(),
	})

	for _, option := range reader.options {
		err := option(reader)
		if err != nil {
			return nil, err
		}
	}

	return reader.blocks, nil
}

func (reader *memoryReader) Write(value interface{}, path ...string) {
	reader.sj.SetPath(path, value)
}

func (reader *memoryReader) Name() string {
	return "memory"
}

// Object .
func Object(object interface{}) Option {
	return func(reader *memoryReader) error {
		data, err := json.Marshal(object)

		if err != nil {
			return err
		}

		reader.blocks = append(reader.blocks, &scf4go.ReadBlock{
			Data:      data,
			Codec:     "json",
			Timestamp: time.Now(),
		})

		return nil
	}
}

// Data .
func Data(data string, codec string) Option {
	return func(reader *memoryReader) error {

		reader.blocks = append(reader.blocks, &scf4go.ReadBlock{
			Data:      []byte(data),
			Codec:     codec,
			Timestamp: time.Now(),
		})

		return nil
	}
}

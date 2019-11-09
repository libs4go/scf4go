package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/libs4go/errors"

	"github.com/libs4go/scf4go"
)

type fileWithCodec struct {
	path  string
	codec string
}

type fileReader struct {
	path    []*fileWithCodec
	options []Option
}

// JSON load json file with path
func JSON(path string) Option {
	return File(path, scf4go.WithCodec("json"))
}

// Yaml load yaml file with path
func Yaml(path string) Option {
	return File(path, scf4go.WithCodec("yaml"))
}

// Option .
type Option func(reader *fileReader) error

// File reader single file
func File(path string, options ...scf4go.Option) Option {

	return func(reader *fileReader) error {

		ext := filepath.Ext(path)

		realcodec := scf4go.NewOptions(options...).Codec

		switch ext {
		case ".yaml":
			realcodec = "yaml"
		}

		reader.path = append(reader.path, &fileWithCodec{
			path:  path,
			codec: realcodec,
		})

		return nil
	}
}

// Dir reader dir files
func Dir(path string, options ...scf4go.Option) Option {

	return func(reader *fileReader) error {

		codec := scf4go.NewOptions(options...).Codec

		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if path == "." || path == ".." {
				return err
			}

			ext := filepath.Ext(path)

			realcodec := codec

			switch ext {
			case ".yaml":
				realcodec = "yaml"
			}

			reader.path = append(reader.path, &fileWithCodec{
				path:  path,
				codec: realcodec,
			})

			return err
		})
	}
}

// New .
func New(options ...Option) scf4go.Reader {
	return &fileReader{options: options}
}

func (reader *fileReader) Read() ([]*scf4go.ReadBlock, error) {

	reader.path = nil
	for _, option := range reader.options {
		if err := option(reader); err != nil {
			return nil, err
		}
	}

	var blocks []*scf4go.ReadBlock

	for _, path := range reader.path {
		data, err := ioutil.ReadFile(path.path)

		if err != nil {
			return nil, errors.Wrap(err, "read file %s error", path)
		}

		blocks = append(blocks, &scf4go.ReadBlock{
			Data:      data,
			Codec:     path.codec,
			Timestamp: time.Now(),
		})
	}

	return blocks, nil
}
func (reader *fileReader) Name() string {
	return "file"
}

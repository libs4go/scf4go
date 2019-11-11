package scf4go

import (
	"encoding/json"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/libs4go/errors"

	"github.com/imdario/mergo"
)

type configImpl struct {
	jsonCodec Codec            // the inner using Codec
	codec     map[string]Codec // support Codec map
	readers   []Reader         // readers
	block     *ReadBlock       // the last merged read block
	sj        *simplejson.Json // parsed simplejson object
}

// New Config create new config object
func New() Config {

	config := &configImpl{
		codec: make(map[string]Codec),
	}

	codecs := Codecs()

	for _, codec := range codecs {
		if codec.Name() == "json" {
			config.jsonCodec = codec
		}
		config.codec[codec.Name()] = codec
	}

	if config.jsonCodec == nil {
		panic(errors.Wrap(ErrJSONCodec, "using register function to register json codec"))
	}

	return config
}

func (config *configImpl) merge(blocks ...*ReadBlock) (*ReadBlock, error) {
	var merged map[string]interface{}

	for _, block := range blocks {
		if block == nil {
			continue
		}

		if len(block.Data) == 0 {
			continue
		}

		codec, ok := config.codec[block.Codec]

		if !ok {
			return nil, errors.Wrap(ErrCodec, "unknown encoder %s", block.Codec)
		}

		var data map[string]interface{}
		if err := codec.Decode(block.Data, &data); err != nil {
			return nil, err
		}

		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	b, err := config.jsonCodec.Encode(merged)
	if err != nil {
		return nil, err
	}

	cs := &ReadBlock{
		Timestamp: time.Now(),
		Data:      b,
		Codec:     "json",
	}

	return cs, nil
}

func (config *configImpl) Close() {

}

func (config *configImpl) Prefix() []string {
	return nil
}

func (config *configImpl) Load(readers ...Reader) error {
	config.readers = append(config.readers, readers...)
	return config.Reload()
}

func (config *configImpl) Reload() error {
	var blocks []*ReadBlock

	for _, reader := range config.readers {
		block, err := reader.Read()

		if err != nil {
			return errors.Wrap(err, "invoke reader %s error", reader.Name())
		}

		blocks = append(blocks, block...)
	}

	block, err := config.merge(blocks...)

	if err != nil {
		return err
	}

	config.block = block
	config.sj = simplejson.New()

	return config.sj.UnmarshalJSON(block.Data)
}

func (config *configImpl) SubConfig(path ...string) Config {
	return newSubConfig(path, config)
}

func (config *configImpl) Get(path ...string) Value {
	return &valueImpl{config.sj.GetPath(path...)}
}

func (config *configImpl) Map() map[string]interface{} {
	m, _ := config.sj.Map()

	return m
}

func (config *configImpl) Scan(v interface{}) error {
	b, err := config.sj.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

type subConfig struct {
	config Config
	prefix []string //  the config search prefix string nodes
}

func newSubConfig(prefix []string, config Config) Config {
	return &subConfig{
		prefix: prefix,
		config: config,
	}
}

func (config *subConfig) Prefix() []string {
	return append(config.config.Prefix(), config.prefix...)
}

func (config *subConfig) Close() {
}

func (config *subConfig) Load(readers ...Reader) error {
	return config.config.Load(readers...)
}

func (config *subConfig) Reload() error {
	return config.config.Reload()
}

func (config *subConfig) SubConfig(path ...string) Config {
	return newSubConfig(path, config)
}

func (config *subConfig) Get(path ...string) Value {

	path = append(config.prefix, path...)
	return config.config.Get(path...)
}

func (config *subConfig) Map() map[string]interface{} {
	var result map[string]interface{}
	err := config.Get(config.Prefix()...).Scan(&result)

	if err != nil {
		panic(errors.Wrap(err, "prefix %s is not a map", strings.Join(config.Prefix(), ".")))
	}

	return result
}

func (config *subConfig) Scan(v interface{}) error {
	return config.Get().Scan(v)
}

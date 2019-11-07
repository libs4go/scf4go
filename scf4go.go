package scf4go

import (
	"sync"
	"time"

	"github.com/libs4go/errors"
	"github.com/libs4go/sdi4go"
)

const errVendor = "scf4go"

// Errors
var (
	ErrCodec     = errors.New("invalid codec name", errors.WithVendor(errVendor), errors.WithCode(-1))
	ErrJSONCodec = errors.New("scf4go basic running mode must import json codec implement", errors.WithVendor(errVendor), errors.WithCode(-1))
)

// Config the config facade
type Config interface {
	Values
	Close()
	Load(readers ...Reader) error
	Reload() error
	SubConfig(path ...string) Config
	Prefix() []string
}

// Values the config values access interface
type Values interface {
	// Retrieve a value
	Get(path ...string) Value
	// Return values as a map
	Map() map[string]interface{}
	// Scan config into a Go type
	Scan(v interface{}) error
}

// Value .
type Value interface {
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
}

// Reader Read the config source object
type Reader interface {
	Read() ([]*ReadBlock, error)
	Name() string
}

// ReadBlock .
type ReadBlock struct {
	Data      []byte
	Codec     string
	Timestamp time.Time
}

// Codec .
type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	Name() string
}

var codecRegister sdi4go.Injector
var initOnce sync.Once

func getRegister() sdi4go.Injector {
	initOnce.Do(func() {
		codecRegister = sdi4go.New()
	})

	return codecRegister
}

// Register .
func Register(codec Codec) {
	getRegister().Bind(codec.Name(), sdi4go.Singleton(codec))
}

// Codecs get register codec slice
func Codecs() []Codec {
	var codecs []Codec

	if err := getRegister().CreateAll(&codecs); err != nil {
		panic(err)
	}

	return codecs
}

// Options .
type Options struct {
	Codec string
}

// Option .
type Option func(options *Options)

// WithCodec .
func WithCodec(name string) Option {
	return func(opts *Options) {
		opts.Codec = name
	}
}

// NewOptions .
func NewOptions(options ...Option) *Options {
	result := &Options{
		Codec: "json",
	}

	for _, opt := range options {
		opt(result)
	}

	return result
}

package yaml

import (
	"github.com/ghodss/yaml"
	"github.com/libs4go/errors"
	"github.com/libs4go/scf4go"
)

type yamlCodec struct{}

func (y yamlCodec) Encode(v interface{}) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err == nil {
		return data, nil
	}

	return nil, errors.Wrap(err, "encode yaml error")
}

func (y yamlCodec) Decode(d []byte, v interface{}) error {
	err := yaml.Unmarshal(d, v)
	if err == nil {
		return nil
	}

	return errors.Wrap(err, "decode yaml error: %s", string(d))
}

func (y yamlCodec) Name() string {
	return "yaml"
}

func init() {
	scf4go.Register(&yamlCodec{})
}

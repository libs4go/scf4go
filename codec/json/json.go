package json

import (
	"encoding/json"

	"github.com/libs4go/errors"

	"github.com/libs4go/scf4go"
)

type jsonCodec struct {
}

func (codec *jsonCodec) Encode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)

	if err != nil {
		return nil, errors.Wrap(err, "encode json error: %v", v)
	}

	return data, nil
}

func (codec *jsonCodec) Decode(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)

	if err != nil {
		return errors.Wrap(err, "decode json error: %s", string(data))
	}

	return nil
}

func (codec *jsonCodec) Name() string {
	return "json"
}

func init() {
	scf4go.Register(&jsonCodec{})
}

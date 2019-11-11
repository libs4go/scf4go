package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/libs4go/scf4go"
	_ "github.com/libs4go/scf4go/codec"
	"github.com/libs4go/scf4go/reader/file"
	"github.com/libs4go/scf4go/reader/memory"
)

func TestMemory(t *testing.T) {
	config := scf4go.New()
	err := config.Load(memory.New(memory.Data(
		`
		{
			"a": {
				"a1": 1,
				"a2": true,
				"a3": "1h"
			},
			"b":12.5
		}
		`,
		"json")))

	require.NoError(t, err)

	require.Equal(t, config.Get("a", "a1").Int(0), 1)
	require.Equal(t, config.Get("a", "a3").Duration(time.Second), time.Hour)
	require.Equal(t, config.Get("a", "a2").Bool(true), true)
	require.Equal(t, config.Get("b").Float64(1.1), 12.5)

	config = config.SubConfig("a")

	require.Equal(t, config.Get("a1").Int(0), 1)
	require.Equal(t, config.Get("a3").Duration(time.Second), time.Hour)
	require.Equal(t, config.Get("a2").Bool(true), true)
}

func TestLoadFile(t *testing.T) {
	config := scf4go.New()
	err := config.Load(file.New(
		file.File("./data/test.json"),
		file.File("./data/test.yaml", scf4go.WithCodec("yaml")),
	))

	require.NoError(t, err)

	require.Equal(t, config.Get("b").Float64(1.1), 13.5)

	servers := config.Get("server").StringSlice([]string{})

	require.Equal(t, len(servers), 3)

	err = config.Load(file.New(file.Yaml("./data/test2.yaml")))

	require.NoError(t, err)

	servers = config.Get("server").StringSlice([]string{})

	require.Equal(t, len(servers), 4)

	config = config.SubConfig("server2").SubConfig("test")

	var r map[string][]string

	err = config.Scan(&r)

	require.NoError(t, err)

	println(fmt.Sprintf("%v", r))
}

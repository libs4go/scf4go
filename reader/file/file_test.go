package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDir(t *testing.T) {
	reader := New(Dir("./data"))
	_, err := reader.Read()
	require.NoError(t, err)
}

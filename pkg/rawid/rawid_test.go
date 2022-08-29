package rawid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBase62_over_range(t *testing.T) {
	t.Parallel()

	input := "ZZZZZZZZZZZ"
	id, err := NewBase62(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "over range")
	assert.Nil(t, id, "returned value should be nil on error")
}

func TestNewBase62_unknown_string(t *testing.T) {
	t.Parallel()

	input := "&"
	id, err := NewBase62(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to decode Base62 input")
	assert.Nil(t, id, "returned value should be nil on error")
}

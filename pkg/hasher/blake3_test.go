package hasher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_blake3_lenOut_in_range(t *testing.T) {
	input := "foo bar"

	{
		r := strings.NewReader(input)

		lenOut := 0
		hashed, err := _blake3(r, lenOut)
		expect := "99a8f992299c9ade63e4285e8193f5a0b9b2ec2c33276ade28894a496" +
			"2ce9c0175abdcb087cee396110dfa8f66c4636ec528915aa95bac7bec2d83c2" +
			"acf25f77"
		actual := fmt.Sprintf("%x", hashed)

		require.NoError(t, err)
		assert.Equal(t, expect, actual)
		assert.Equal(t, 64, len(hashed), "if lenOut is zero then the output length should be 64 bytes length")
	}
	{
		r := strings.NewReader(input)

		lenOut := 8192
		hashed, err := _blake3(r, lenOut)
		expectHead := "99a8f992299c9ade63e4285e8193f5a0"
		expectTail := "0bed4746ce6e4ee34c8b5c09c4ba91bf"
		actual := fmt.Sprintf("%x", hashed)

		require.NoError(t, err)
		assert.Equal(t, expectHead, actual[:32])
		assert.Equal(t, expectTail, actual[len(actual)-32:])
	}
}

func Test_blake3_lenOut_out_of_range(t *testing.T) {
	input := "foo bar"

	{
		r := strings.NewReader(input)

		lenOut := -1
		hashed, err := _blake3(r, lenOut)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid output length. It must be between 1 and 8194")
		assert.Nil(t, hashed)
	}
	{
		r := strings.NewReader(input)

		lenOut := 8194 + 1
		hashed, err := _blake3(r, lenOut)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid output length. It must be between 1 and 8194")
		assert.Nil(t, hashed)
	}
}

func Test_blake3_scan_error(t *testing.T) {
	// See hasher_test.go for dummyReader struct
	d := dummyReader{}

	hashed, err := _blake3(d, 16)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read during scanning")
	assert.Nil(t, hashed)
}

func Test_blake3_seek_error(t *testing.T) {
	// Backup and defer restore before mocking io.SeekStart
	oldIoSeekStart := IoSeekStart
	defer func() {
		IoSeekStart = oldIoSeekStart
	}()

	// Mock io.SeekStart
	IoSeekStart = 10

	input := "foo bar"
	r := strings.NewReader(input)

	hashed, err := _blake3(r, 16)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set the position to seek")
	assert.Nil(t, hashed)
}

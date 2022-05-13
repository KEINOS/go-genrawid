package hasher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_crc32_golen(t *testing.T) {
	input := "Hello world!"
	r := strings.NewReader(input)

	chksum, err := _crc32(r)

	require.NoError(t, err)
	assert.Equal(t, 4, len(chksum))

	// Expect CRC32C value taken from Rust implementation:
	//   https://docs.rs/crc32c/latest/crc32c/
	expect := "7b98e751"
	actual := fmt.Sprintf("%x", chksum)
	assert.Equal(t, expect, actual)
}

func Test_crc32_fail_copy(t *testing.T) {
	// See hasher_test.go for dummyReader struct
	d := dummyReader{}
	checksum, err := _crc32(d)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to copy data to hasher")
	assert.Nil(t, checksum)
}

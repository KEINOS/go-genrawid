package hasher

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  CheckSum
// ----------------------------------------------------------------------------

func TestCheckSum_default_golden(t *testing.T) {
	t.Parallel()

	//nolint:varnamelen // allow short variable names for readability
	for i, test := range []struct {
		input  string
		expect string
	}{
		{"1234567890", "f3dbd4fe"},
		{"123456789", "e3069283"},
		{"The quick brown fox jumps over the lazy dog", "22620404"},
	} {
		input := test.input
		r := strings.NewReader(input)

		sumByte, err := CheckSum(r)
		if err != nil {
			log.Fatal(err)
		}

		expect := test.expect
		actual := fmt.Sprintf("%x", sumByte)
		assert.Equal(t, expect, actual, "case #%d: input: %s", i, input)
	}
}

func TestCheckSum_nil_input(t *testing.T) {
	t.Parallel()

	checksum, err := CheckSum(nil)

	require.Error(t, err, "nil input should be an error")
	assert.Contains(t, err.Error(), "nil pointer for input given")
	assert.Nil(t, checksum, "returned checksum should be nil on error")
}

//nolint:paralleltest // do not parallelize due to dependency on other tests
func TestCheckSum_unknown_algo(t *testing.T) {
	oldChkSumAlgo := ChkSumAlgo
	defer func() {
		ChkSumAlgo = oldChkSumAlgo
	}()

	// Set unknown dummy algorithm
	ChkSumAlgo = ChkSumUnknown

	checksum, err := CheckSum(strings.NewReader("This is a string"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown checksum algorithm")
	assert.Nil(t, checksum, "returned checksum should be nil on error")
}

// ----------------------------------------------------------------------------
//  Hash
// ----------------------------------------------------------------------------

func TestHash_default_golden(t *testing.T) {
	t.Parallel()

	input := "This is a string"
	hashByte, err := Hash(strings.NewReader(input))

	require.NoError(t, err)

	expect := "718b749f12a61257438b2ea6643555fd995001c9d9ff84764f93f82610a78" +
		"0f243a9903464658159cf8b216e79006e12ef3568851423fa7c97002cbb9ca4dc44"
	actual := fmt.Sprintf("%x", hashByte)
	assert.Equal(t, expect, actual, "input string: %s", input)
}

func TestHash_nil_input(t *testing.T) {
	t.Parallel()

	checksum, err := Hash(nil)

	require.Error(t, err, "nil input should be an error")
	assert.Contains(t, err.Error(), "nil pointer for input given")
	assert.Nil(t, checksum, "returned checksum should be nil on error")
}

//nolint:paralleltest // do not parallelize due to dependency on other tests
func TestHash_unknown_algo(t *testing.T) {
	oldHashAlgo := HashAlgo
	defer func() {
		HashAlgo = oldHashAlgo
	}()

	// Set unknown dummy algorithm
	HashAlgo = HashAlgoUnknown

	digest, err := Hash(strings.NewReader("This is a string"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown hash algorithm")
	assert.Nil(t, digest, "returned digest should be nil on error")
}

// ----------------------------------------------------------------------------
//  Helpers (Dummy reader struct)
// ----------------------------------------------------------------------------

type dummyReader struct{}

//nolint:nonamedreturns // allow named return for interface compatibility
func (r dummyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

//nolint:nonamedreturns // allow named return for interface compatibility
func (r dummyReader) WriteTo(w io.Writer) (n int64, err error) {
	return 0, errors.New("forced error")
}

type dummyReader2 struct{}

//nolint:nonamedreturns // allow named return for interface compatibility
func (r dummyReader2) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced error")
}

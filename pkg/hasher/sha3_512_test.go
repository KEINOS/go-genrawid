package hasher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Tests for private functions
// ----------------------------------------------------------------------------

func Test_sha3_512_golden(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{
			// Expect value taken from Wikipeadia.
			// https://en.wikipedia.org/wiki/SHA-3#Examples_of_SHA-3_variants
			input: "",
			expect: "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc" +
				"1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301" +
				"758586281dcd26",
		},
		{
			// Expect value taken from OpenSSL v3.0.3:
			//   echo -n 'beef' | openssl sha3-512
			input: "beef",
			expect: "b80b65157d02870f4691e10490ac1b0f5ed280523d3a1dfa2211c6a" +
				"d241a825081e5de9f226b12943f8574cac0694aad020e7816f911bdccbc" +
				"c73c1981d18c12",
		},
	} {
		hashByte, err := _sha3_512(strings.NewReader(test.input), 0)

		require.NoError(t, err)
		assert.Equal(t, 64, len(hashByte), "if lenOut is 0 it should treat as 64 by default")

		expect := test.expect
		actual := fmt.Sprintf("%x", hashByte)
		assert.Equal(t, expect, actual, "input string: %s", test.input)
	}
}

func Test_sha3_512_output_length_out_of_range(t *testing.T) {
	input := "foo bar"

	{
		lenBad := -1 // output length should be between 0-64.

		hashByte, err := _sha3_512(strings.NewReader(input), lenBad)

		require.Error(t, err, "negative length should be an error")
		assert.Contains(t, err.Error(), "invalid output length. It must be between 1 and 64")
		assert.Nil(t, hashByte, "returned hash should be nil on error")
	}
	{
		lenBad := 65 // output length should be between 0-64.

		hashByte, err := _sha3_512(strings.NewReader(input), lenBad)

		require.Error(t, err, "length over 64 should be an error")
		assert.Contains(t, err.Error(), "invalid output length. It must be between 1 and 64")
		assert.Nil(t, hashByte, "returned hash should be nil on error")
	}
}

func Test_sha3_512_read_error(t *testing.T) {
	// See hasher_test.go for dummyReader struct
	d := dummyReader{}

	hashByte, err := _sha3_512(d, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read during scanning")
	assert.Nil(t, hashByte, "return value should be nil on error")
}

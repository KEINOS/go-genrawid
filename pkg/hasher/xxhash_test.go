package hasher

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example of the use of the private function _xxhash. To use it from public see
// the CheckSum (Xxhash) example.
func Example_xxhash() {
	input := "a"
	r := strings.NewReader(input)

	// lenOut=0 returns 64 bit/8 Byte length output
	byteR, err := _xxhash(r, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("0x%x\n", byteR)

	// Output:
	// 0xd24ec4f1a98c6e5b
}

func Test_xxhash_golden(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		inputStr  string
		expectHex string
		inputLen  int
		expectLen int
	}{
		{"", "0xef46db3751d8e999", 0, 8},
		{"", "0xef", 1, 1},
		{"", "0xef46db37", 4, 4},
		{"", "0xef46db3751d8e999", 8, 8},
		{"a", "0xd24ec4f1a98c6e5b", 0, 8},
		{"as", "0x1c330fb2d66be179", 0, 8},
		{"asd", "0x631c37ce72a97393", 0, 8},
		{"asdf", "0x415872f599cea71e", 0, 8},
		{"1234567890", "0xa9d4d4132eff23b6", 0, 8},
	} {
		r := strings.NewReader(test.inputStr)

		byteR, err := _xxhash(r, test.inputLen)
		require.NoError(t, err)

		expect := test.expectHex
		actual := fmt.Sprintf("0x%x", byteR)
		require.Equal(t, expect, actual, "input string was:", test.inputStr)

		assert.Equal(t, test.expectLen, len(byteR), "unexpected output length")
	}
}

func Test_xxhash_lenOut_too_long(t *testing.T) {
	t.Parallel()

	input := "a"
	r := strings.NewReader(input)

	byteR, err := _xxhash(r, 9)

	require.Error(t, err)
	assert.Nil(t, byteR, "returned value should be nil on error")
}

func Test_xxhash_read_error(t *testing.T) {
	t.Parallel()

	// See hasher_test.go for dummyReader2 struct
	d := dummyReader2{}

	hashByte, err := _xxhash(d, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error during reading data")
	assert.Nil(t, hashByte, "return value should be nil on error")
}

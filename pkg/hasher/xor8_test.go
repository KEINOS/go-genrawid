package hasher

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example_xor8() {
	input := "this is a sample data"
	r := strings.NewReader(input)

	byteR, err := _xor8(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", byteR)
	fmt.Printf("%08b\n", byteR)
	fmt.Printf("%04d\n", byteR)

	// Output:
	// 6f
	// [01101111]
	// [0111]
}

func Test_xor8_read_error(t *testing.T) {
	// See hasher_test.go for dummyReader2 struct
	d := dummyReader2{}

	hashByte, err := _xor8(d)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error during reading data")
	assert.Nil(t, hashByte, "return value should be nil on error")
}

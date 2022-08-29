//go:generate go run ./testdata/gen/genFileDummy.go ./testdata/dummy.bin
package genrawid

import (
	"os"
	"strings"
	"testing"

	"github.com/KEINOS/go-genrawid/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Tests for public function
// ----------------------------------------------------------------------------

func TestFromFile_big_file(t *testing.T) {
	t.Parallel()

	pathFile := "testdata/dummy.bin"

	rawid, err := FromFile(pathFile)
	require.NoError(t, err)

	assert.Equal(t, "-2929669798473946006", rawid.Dec())
}

//nolint:paralleltest // do not parallelize due to dependency on other tests
func TestFromStdin_big_file(t *testing.T) {
	// Backup and defer recover the stdin
	oldStdin := OsStdin
	defer func() {
		OsStdin = oldStdin
	}()

	// Mock the stdin by assigning a file to stdin.
	pathFile := "testdata/dummy.bin"

	osFile, err := os.Open(pathFile)
	require.NoError(t, err)

	OsStdin = osFile

	rawid, err := FromStdin()
	require.NoError(t, err)

	assert.Equal(t, "-2929669798473946006", rawid.Dec())
}

func TestFromFile_file_not_found(t *testing.T) {
	t.Parallel()

	pathFile := "dummy/unknown/file"

	rawid, err := FromFile(pathFile)
	require.Error(t, err, "it should be an error on file not found")

	assert.Contains(t, err.Error(), "failed to open file")
	assert.Empty(t, rawid, "rawid should be empty on error")
}

// ----------------------------------------------------------------------------
//  Tests for private function
// ----------------------------------------------------------------------------

func Test_chopAndMergeBytes_golden(t *testing.T) {
	t.Parallel()

	a := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8}
	b := []byte{0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF, 0x10}

	// Merge the first 4 bytes of a and b.
	rawid, err := chopAndMergeBytes(a, b)

	require.NoError(t, err)

	expect := []byte{0x1, 0x2, 0x3, 0x4, 0x9, 0xA, 0xB, 0xC}
	actual := rawid.Byte()
	assert.Equal(t, expect, actual)
}

func Test_chopAndMergeBytes_too_few_slice(t *testing.T) {
	t.Parallel()

	{
		a := []byte{0x1, 0x2, 0x3}
		b := []byte{0x9, 0xA, 0xB, 0xC, 0xD, 0xE}

		rawid, err := chopAndMergeBytes(a, b)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to combine bytes. Both of the input must be 4byte or more")
		assert.Nil(t, rawid, "on error the returned rawid should be nil")
	}
	{
		a := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
		b := []byte{0x9, 0xA, 0xB}

		rawid, err := chopAndMergeBytes(a, b)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to combine bytes. Both of the input must be 4byte or more")
		assert.Nil(t, rawid, "on error the returned rawid should be nil")
	}
}

func Test_genRawid_nil_input(t *testing.T) {
	t.Parallel()

	rawid, err := genRawid(nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil pointer for input given")
	assert.Nil(t, rawid)
}

//nolint:paralleltest // do not parallelize due to dependency on other tests
func Test_genRawid_use_undefined_algo(t *testing.T) {
	// Set unknown checksum algorithm
	oldChkSumAlgo := hasher.ChkSumAlgo
	defer func() {
		hasher.ChkSumAlgo = oldChkSumAlgo
	}()

	hasher.ChkSumAlgo = hasher.ChkSumUnknown

	// Dummy input
	r := strings.NewReader("sample input")

	// Test
	rawid, err := genRawid(r)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate rawid")
	assert.Nil(t, rawid)
}

func Test_replaceLast16bit(t *testing.T) {
	t.Parallel()

	inHash := []byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	inChkSum := uint16(0xffff)

	expect := []byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0xff, 0xff}
	actual := replaceLast16bit(inHash, inChkSum)

	assert.Equal(t, expect, actual, "the last 2 bytes of the rawid should be replaced with the checksum")
}

func Test_xorSliceByte(t *testing.T) {
	t.Parallel()

	input := []byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05}

	expect := uint16(0x0707) // 0x01 ^ 0x02 ^ 0x4 = 0x07, 0x01 ^ 0x03 ^ 0x5 = 0x07
	actual := xorSliceByte(input)

	assert.Equal(t, expect, actual)
}

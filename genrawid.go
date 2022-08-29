/*
Package genrawid provides functions to generate rawid.

For the sample implementation see: ./cmd/genrawid/main.go
*/
package genrawid

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/KEINOS/go-genrawid/pkg/hasher"
	"github.com/KEINOS/go-genrawid/pkg/rawid"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Variables
// ----------------------------------------------------------------------------

// OsStdin is a copy of os.Stdin to ease testing. Mock this variable during tests.
var OsStdin = os.Stdin

// IsModeFast is the flag to use fast mode. It will use the last 16bit of the hash
// as the xor16 checksum.
var IsModeFast = false

// ----------------------------------------------------------------------------
//  Functions (Public)
// ----------------------------------------------------------------------------

// FromFile returns the rawid generated from the input file.
func FromFile(path string) (rawid.ID, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer file.Close()

	return genRawid(file)
}

// FromStdin returns the rawid generated from stdin as its input.
func FromStdin() (rawid.ID, error) {
	return genRawid(OsStdin)
}

// FromString returns the rawid generated from the input string.
func FromString(input string) (rawid.ID, error) {
	r := strings.NewReader(input)

	return genRawid(r)
}

// ----------------------------------------------------------------------------
//  Functions (Private)
// ----------------------------------------------------------------------------

// It combines the two input as one in 8 byte length.
//
//nolint:varnamelen // allow short variable names for readability.
func chopAndMergeBytes(a, b []byte) (rawid.ID, error) {
	if len(a) < 4 || len(b) < 4 {
		return nil, errors.New("failed to combine bytes. Both of the input must be 4byte or more")
	}

	lenByte := 8
	rawid := make([]byte, lenByte)

	copy(rawid, a)     // Upper half as hash
	copy(rawid[4:], b) // Bottom half as checksum

	return rawid, nil
}

// It returns 64bit/8byte length rawid.
func genRawid(input io.Reader) (rawid.ID, error) {
	// Calculate hash value.
	hashByte, err := hasher.Hash(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate rawid")
	}

	// Calculate checksum of the hash.
	if !IsModeFast {
		r := bytes.NewReader(hashByte)

		sumByte, err := hasher.CheckSum(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate rawid")
		}

		// Combine the fisrt 4 Bytes of the hash and the checksum as a rawid.
		return chopAndMergeBytes(hashByte, sumByte)
	}

	// ------------------------------------------------------------------------
	// Fast mode: use the last 16 bit/2 Bytes of the hash as the xor16 checksum.
	// Important:
	//   This mode is currently hidden since it was not fast enough as expected.
	// ------------------------------------------------------------------------

	// Calculate the xor16 checksum of the hash.
	chkSum := xorSliceByte(hashByte)
	lenByte := 8
	rawid := make([]byte, lenByte)

	// Set hash
	copy(rawid, hashByte)

	// Set the last 2 bytes(16bit) of the hash as the checksum.
	return replaceLast16bit(rawid, chkSum), nil
}

func replaceLast16bit(input []byte, xor16 uint16) []byte {
	lenBitShift := 8

	input[len(input)-2] = byte(xor16 >> lenBitShift)
	input[len(input)-1] = byte(xor16)

	return input
}

func xorSliceByte(input []byte) uint16 {
	lenBitShift := 8

	var out uint16 = 0

	for i, b := range input {
		if i%2 == 0 {
			out ^= uint16(b)
		} else {
			out ^= uint16(b) << lenBitShift
		}
	}

	return out
}

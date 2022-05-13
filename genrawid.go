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
func chopAndMergeBytes(a, b []byte) (rawid.ID, error) {
	if len(a) < 4 || len(b) < 4 {
		return nil, errors.New("failed to combine bytes. Both of the input must be 4byte or more")
	}

	rawid := make([]byte, 8)

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
	r := bytes.NewReader(hashByte)

	sumByte, err := hasher.CheckSum(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate rawid")
	}

	// Combine the fisrt 4 Bytes of the hash and the checksum as a rawid.
	return chopAndMergeBytes(hashByte, sumByte)
}

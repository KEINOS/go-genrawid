/*
Package hasher implements hash functions. It is used to generate raw IDs.
*/
package hasher

import (
	"hash/crc32"
	"io"

	"github.com/KEINOS/go-genrawid/pkg/rawid"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Unexposed Constants
// ----------------------------------------------------------------------------

// ChksumAlgoDefault is the default checksum algorithm. Currently, it is CRC32.
const chksumAlgoDefault = ChkSumCRC32

// CRC32PolyDefault is the default polynomial used in the CRC32 algorithm, which
// is the Castagnoli polynomial.
//
// Castagnoli polynomial, like Koopman polynomial, can detect more characteristic
// errors than IEEE polynomial.
const crc32PolyDefault = uint32(crc32.Castagnoli)

// HashAlgoTypeDefault is the default hash algorithm. Currently it is BLAKE3.
const hashAlgoDefault = HashAlgoBLAKE3

// HashDigestSize is the default byte length of the hash digest.
const hashLenDefault = 64

// ----------------------------------------------------------------------------
//  Exposed Variables with default values set
// ----------------------------------------------------------------------------
//  These variables can be set by the user to change the default values.

// ChkSumAlgo is the checksum algorithm to use.
//
// One of TChkSumAlgo type must be set. By default it is ChkSumCRC32(=CRC32).
// Currently, we only support CRC32.
var ChkSumAlgo = chksumAlgoDefault

// CRC32Poly is the polynomial used in the CRC32 algorithm.
//
// By default it uses Castagnoli polynomial. Overwrite this value to use a
// different polynomial.
//
//   Example:
//     genrawid.CRC32Poly = crc32.Koopman
//     genrawid.CRC32Poly = crc32.IEEE
//
// Note that this won't affect if ChkSumAlgo is other than CRC32.
var CRC32Poly = crc32PolyDefault

// HashAlgo is the hash algorithm to use. By default it uses BLAKE3.  Overwrite
// this value to use a different hash algorithm.
var HashAlgo = hashAlgoDefault

// HashLen is the hash digest length. By default it is 64 byte length.
var HashLen = hashLenDefault

// ----------------------------------------------------------------------------
//  Public Functions
// ----------------------------------------------------------------------------

// Hash returns the hash/digest of input. By default the returned digest length
// is 64 bytes.
func Hash(input io.Reader) (rawid.ID, error) {
	if input == nil {
		return nil, errors.New("nil pointer for input given")
	}

	switch HashAlgo {
	case HashAlgoBLAKE3:
		return _blake3(input, HashLen)
	case HashAlgoSHA3_512:
		return _sha3_512(input, HashLen)
	}

	return nil, errors.Errorf("unknown hash algorithm: %s", HashAlgo)
}

// CheckSum returns the CRC-32 checksum of input.
//
// The CRC32Poly variable is used as a polynomial to create the table. By default
// it uses Castagnoli polynomial. A.k.a. CRC32C or CRC32-Castagnoli.
func CheckSum(input io.Reader) (rawid.ID, error) {
	const lenByte = 4 // output length

	if input == nil {
		return nil, errors.New("nil pointer for input given")
	}

	switch ChkSumAlgo {
	case ChkSumCRC32:
		return _crc32(input)
	case ChkSumXXHash:
		return _xxhash(input, lenByte)
	}

	return nil, errors.Errorf("unknown checksum algorithm: %s", ChkSumAlgo)
}

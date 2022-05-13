package hasher

// ----------------------------------------------------------------------------
//  Types
// ----------------------------------------------------------------------------

// THashAlgo is an enum type that represents the hash algorithm.
type THashAlgo int

// TChkSumAlgo is an enum type that represents the checksum algorithm.
type TChkSumAlgo int

// ----------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------

const (
	// HashAlgoUnknown is the enum of unknown hash algorithm.
	HashAlgoUnknown THashAlgo = iota
	// HashAlgoBLAKE3 is the enum of BLAKE3 hash algorithm. This is the default
	// set to "hasher.HashAlgo".
	HashAlgoBLAKE3
	// HashAlgoSHA3_512 is the enum of SHA3-512 hash algorithm. Set this to
	// "hasher.HashAlgo" to use this algorithm.
	HashAlgoSHA3_512
)

const (
	// ChkSumUnknown is the enum of unknown checksum algorithm.
	ChkSumUnknown TChkSumAlgo = iota
	// ChkSumCRC32 is the enum of CRC32 checksum algorithm. This needs "CRC32Poly"
	// variagle, a polynomial to use, to be set as well.
	ChkSumCRC32
	// ChkSumXXHash is the enum of xxHash algorithm as a checksum.
	ChkSumXXHash
)

// ----------------------------------------------------------------------------
//  Methods (Implementation of fmt.Stringer)
// ----------------------------------------------------------------------------

// String returns the string representation of the hash algorithm.
func (h THashAlgo) String() string {
	switch h {
	case HashAlgoBLAKE3:
		return "blake3"
	case HashAlgoSHA3_512:
		return "sha3-512"
	}

	return "unknown"
}

// String returns the string representation of the hash algorithm.
func (h TChkSumAlgo) String() string {
	switch h {
	case ChkSumCRC32:
		return "crc32"
	case ChkSumXXHash:
		return "xxhash"
	}

	return "unknown"
}

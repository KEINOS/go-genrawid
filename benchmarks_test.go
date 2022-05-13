package genrawid_test

// ============================================================================
//  Mislenious Benchmarks to ensure which method/function is faster.
// ============================================================================

import (
	"testing"

	"github.com/KEINOS/go-genrawid"
	"github.com/KEINOS/go-genrawid/pkg/hasher"
)

// ----------------------------------------------------------------------------
//  Simple use between & and %
// ----------------------------------------------------------------------------
//  Current conclusion: Not much difference between & and %
//  This is because the compiler automatically converts to "&" internally if
//  possible. See "Use & and % for managing values" bench for the differences.
//
//    name        time/op
//    AND_simple  0.33ns ± 1%
//    MOD_simple  0.33ns ± 3%

func BenchmarkAND_simple(b *testing.B) {
	ite := 3 // (2^n)-1, n=2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = i & ite
	}
}

func BenchmarkMOD_simple(b *testing.B) {
	ite := 4 // 2^n, n=2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = i % ite
	}
}

// ----------------------------------------------------------------------------
//  Use & and % for managing values
// ----------------------------------------------------------------------------
//  Current conclusion: & is 93% faster than %
//
//    name        time/op
//    AND_manage  0.74ns ± 1%
//    MOD_manage  11.0ns ± 2%

func BenchmarkAND_manage(b *testing.B) {
	ite := 3 // (2^n)-1, n=2

	var d [4]int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:govet // we do not use the variable during bench
		d[i&ite] &= (i & ite) + 1
	}
}

func BenchmarkMOD_manage(b *testing.B) {
	ite := 4 // 2^n, n=2

	var d [4]int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//nolint:govet // we do not use the variable during bench
		d[i%ite] %= (i % ite) + 1
	}
}

// ----------------------------------------------------------------------------
//  FromString
// ----------------------------------------------------------------------------

// Hash: BLAKE3-512, Checksum: CRC32C
// Log:
//   date: Sat May  7 03:45:45 UTC 2022
//   goos: linux
//   goarch: amd64
//   pkg: github.com/KEINOS/go-genrawid
//   cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//   BenchmarkFromString_BLAKE3_512-2   	    1912	    621682 ns/op	    4272 B/op	       7 allocs/op
func BenchmarkFromString_BLAKE3_512(b *testing.B) {
	// Ensure to use BLAKE3 hash algorithm
	oldHashAlgo := hasher.HashAlgo
	defer func() {
		hasher.HashAlgo = oldHashAlgo
	}()

	hasher.HashAlgo = hasher.HashAlgoBLAKE3 // set BLAKE3

	if hasher.ChkSumAlgo != hasher.ChkSumCRC32 {
		b.Fatal("Checksum is not CRC32C")
	}

	input := string(testData(b)) // 1MB of data

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = genrawid.FromString(input)
	}
}

// Hash: SHA3-512, Checksum: CRC32C
// Log:
//   date: Sat May  7 03:52:43 UTC 2022
//   goos: linux
//   goarch: amd64
//   pkg: github.com/KEINOS/go-genrawid
//   cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//   BenchmarkFromString_SHA3_512-2   	     158	   7497839 ns/op	    5232 B/op	      10 allocs/op
func BenchmarkFromString_SHA3_512(b *testing.B) {
	// Switch hash algorithm to SHA3
	oldHashAlgo := hasher.HashAlgo
	defer func() {
		hasher.HashAlgo = oldHashAlgo
	}()

	hasher.HashAlgo = hasher.HashAlgoSHA3_512 // set SHA3

	if hasher.ChkSumAlgo != hasher.ChkSumCRC32 {
		b.Fatal("Checksum is not CRC32C")
	}

	input := string(testData(b)) // 1MB of data

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = genrawid.FromString(input)
	}
}

// Hash: BLAKE3, Checksum: xxHash
// Log:
//   date: Tue May 10 04:43:54 UTC 2022
//   goos: linux
//   goarch: amd64
//   pkg: github.com/KEINOS/go-genrawid
//   cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//   BenchmarkFromString_xxhash_as_checksum-2   	    1702	    687642 ns/op	    4264 B/op	       7 allocs/op
func BenchmarkFromString_xxhash_as_checksum(b *testing.B) {
	// Switch checksum from CRC32C to xxHash
	oldChkSumAlgo := hasher.ChkSumAlgo
	defer func() {
		hasher.ChkSumAlgo = oldChkSumAlgo
	}()

	hasher.ChkSumAlgo = hasher.ChkSumXXHash // set xxHash as checksum

	if hasher.HashAlgo != hasher.HashAlgoBLAKE3 {
		b.Fatal("Hash function is not BLAKE3")
	}

	input := string(testData(b)) // 1MB of data

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = genrawid.FromString(input)
	}
}

// ============================================================================
//  Helper Functions
// ============================================================================

// inputData holds 1MB(1e6) size of data created by testData() function.
var inputData []byte

// The testData creates 1,000,000 bytes= 1MB (1e6) size of data.
// The returned values are consistent and not random.
func testData(b *testing.B) []byte {
	b.Helper()

	// use initialized data
	if len(inputData) != 0 {
		return inputData
	}

	// Initialize data
	inputData = make([]byte, 1e6)

	for i := range inputData {
		// Custom this line to generate different data
		inputData[i] = byte(i % 251)
	}

	return inputData
}

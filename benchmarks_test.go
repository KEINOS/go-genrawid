package genrawid_test

// ============================================================================
//  Mislenious Benchmarks to ensure which method/function is faster.
// ============================================================================

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/KEINOS/go-genrawid"
	"github.com/KEINOS/go-genrawid/pkg/hasher"
	"github.com/stretchr/testify/require"
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
		d[i&ite] &= (i & ite) + 1
	}
}

func BenchmarkMOD_manage(b *testing.B) {
	ite := 4 // 2^n, n=2

	var d [4]int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d[i%ite] %= (i % ite) + 1
	}
}

// ----------------------------------------------------------------------------
//  Fast Mode
// ----------------------------------------------------------------------------

// Benchresults:
//
//	Benchmark_mode_fast-2   	    1820	    617081 ns/op	    4200 B/op	       4 allocs/op
//	Benchmark_mode_fast-2   	    1326	   1088111 ns/op	    4200 B/op	       4 allocs/op
//	Benchmark_mode_fast-2   	    1748	    669352 ns/op	    4200 B/op	       4 allocs/op
func Benchmark_mode_fast(b *testing.B) {
	oldMode := genrawid.IsModeFast
	defer func() {
		genrawid.IsModeFast = oldMode
	}()

	// Ensure fast mode
	genrawid.IsModeFast = true

	input := string(testData(b)) // 1MB of data

	id, err := genrawid.FromString(input)
	require.NoError(b, err)
	require.Equal(b, "012346b77bc3f302", id.Hex())

	// Begin benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = genrawid.FromString(input)
	}
}

// Benchresults:
//
//	Benchmark_mode_regular-2   	    2032	    603058 ns/op	    4272 B/op	       7 allocs/op
//	Benchmark_mode_regular-2   	    2011	    709988 ns/op	    4272 B/op	       7 allocs/op
//	Benchmark_mode_regular-2   	    1970	    636103 ns/op	    4272 B/op	       7 allocs/op
func Benchmark_mode_regular(b *testing.B) {
	oldMode := genrawid.IsModeFast
	defer func() {
		genrawid.IsModeFast = oldMode
	}()

	// Ensure regular mode
	genrawid.IsModeFast = false

	input := string(testData(b)) // 1MB of data

	//nolint:varnamelen // allow short variable names for readability
	id, err := genrawid.FromString(input)
	if err != nil {
		b.Fatal(err)
	}

	if id.Hex() != "012346b7dd5f4297" {
		b.Fatalf("Expected 012346b7dd5f4297, got %s", id.Hex())
	}

	// Begin benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = genrawid.FromString(input)
	}
}

// ----------------------------------------------------------------------------
//  FromString
// ----------------------------------------------------------------------------

// Hash: BLAKE3-512, Checksum: CRC32C
// Log:
//
//	date: Sat May  7 03:45:45 UTC 2022
//	goos: linux
//	goarch: amd64
//	pkg: github.com/KEINOS/go-genrawid
//	cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//	BenchmarkFromString_BLAKE3_512-2   	    1912	    621682 ns/op	    4272 B/op	       7 allocs/op
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
//
//	date: Sat May  7 03:52:43 UTC 2022
//	goos: linux
//	goarch: amd64
//	pkg: github.com/KEINOS/go-genrawid
//	cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//	BenchmarkFromString_SHA3_512-2   	     158	   7497839 ns/op	    5232 B/op	      10 allocs/op
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
//
//	date: Tue May 10 04:43:54 UTC 2022
//	goos: linux
//	goarch: amd64
//	pkg: github.com/KEINOS/go-genrawid
//	cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//	BenchmarkFromString_xxhash_as_checksum-2   	    1702	    687642 ns/op	    4264 B/op	       7 allocs/op
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

// ----------------------------------------------------------------------------
//  xorSliceByte
// ----------------------------------------------------------------------------

// Benchresults:
//
//	Benchmark_uint16-4   	    3349	    362437 ns/op	       0 B/op	       0 allocs/op
//	Benchmark_uint16-4   	    3526	    337835 ns/op	       0 B/op	       0 allocs/op
func Benchmark_uint16(b *testing.B) {
	testFunc := func(input []byte) uint16 {
		var out uint16

		for i, b := range input {
			if i%2 == 0 {
				out ^= uint16(b)
			} else {
				out ^= uint16(b) << 8
			}
		}

		return out
	}

	// Test before bench

	//nolint:ifshort // leave as is for readability
	expect := uint16(0x0707) // 0x01 ^ 0x02 ^ 0x4 = 0x07, 0x01 ^ 0x03 ^ 0x5 = 0x07
	actual := testFunc([]byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05})

	if expect != actual {
		b.Fatal("Expect:", expect, "Actual:", actual)
	}

	// Benchmark
	input := testData(b) // Get test data

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testFunc(input)
	}
}

// Benchresults:
//
//	Benchmark_bigNewInt_SetBytes-4   	    1024	   1163578 ns/op	       8 B/op	       1 allocs/op
//	Benchmark_bigNewInt_SetBytes-4   	     997	   1170646 ns/op	       8 B/op	       1 allocs/op
func Benchmark_bigNewInt_SetBytes(b *testing.B) {
	testFunc := func(input []byte) uint16 {
		out := make([]byte, 2)

		for i, b := range input {
			index := i % 2
			out[index] ^= b
		}

		return uint16(big.NewInt(0).SetBytes(out).Int64())
	}

	// Test before bench
	expect := uint16(0x0707) // 0x01 ^ 0x02 ^ 0x4 = 0x07, 0x01 ^ 0x03 ^ 0x5 = 0x07
	actual := testFunc([]byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05})
	require.Equal(b, expect, actual)

	// Benchmark
	input := testData(b) // Get test data

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testFunc(input)
	}
}

// ----------------------------------------------------------------------------
//  uint16ToBytes
// ----------------------------------------------------------------------------

// Benchresults:
//
//	Benchmark_uint16ToBytes_PutUint16-4   	1000000000	         0.3451 ns/op	       0 B/op	       0 allocs/op
//	Benchmark_uint16ToBytes_PutUint16-4   	1000000000	         0.3514 ns/op	       0 B/op	       0 allocs/op
func Benchmark_uint16ToBytes_PutUint16(b *testing.B) {
	testFunc := func(input uint16) []byte {
		out := make([]byte, 2)
		binary.LittleEndian.PutUint16(out, input)

		return out
	}

	expect := []byte{0xab, 0xcd}
	actual := testFunc(uint16(0xabcd))

	require.Equal(b, expect, actual)

	// Benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testFunc(0xabcd)
	}
}

// Benchresults:
//
//	Benchmark_uint16ToBytes_shift-4   	1000000000	         0.3353 ns/op	       0 B/op	       0 allocs/op
//	Benchmark_uint16ToBytes_shift-4   	1000000000	         0.3327 ns/op	       0 B/op	       0 allocs/op
func Benchmark_uint16ToBytes_shift(b *testing.B) {
	testFunc := func(input uint16) []byte {
		out := make([]byte, 2)
		out[0], out[1] = uint8(input>>8), uint8(input&0xff)

		return out
	}

	expect := []byte{0xab, 0xcd}
	actual := testFunc(uint16(0xabcd))

	require.NotEqual(b, expect, actual)

	// if bytes.Compare(expect, actual) > 0 {
	// 	b.Fatal("Expect:", expect, "Actual:", actual)
	// }

	// Benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testFunc(0xabcd)
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

package hasher_test

// ============================================================================
//  Mislenious Benchmarks to ensure which method/function is faster.
// ============================================================================

import (
	"strings"
	"testing"

	"github.com/KEINOS/go-genrawid/pkg/hasher"
)

// ----------------------------------------------------------------------------
//  Comparison between BLAKE3 and SHA3-512
// ----------------------------------------------------------------------------
//  Current conclusion: Use BLAKE3 as the default algorithm.
//  According to the below statistics, BLAKE3 was twice as fast than SHA3-512.
//
//    goos: darwin
//    goarch: amd64
//    cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
//
//    name          time/op
//    HashBLAKE3-4  1.54µs ±11%
//    HashSHA3-4    2.11µs ± 1%
//
//    name          alloc/op
//    HashBLAKE3-4  4.19kB ± 0%
//    HashSHA3-4    5.15kB ± 0%
//
//    name          allocs/op
//    HashBLAKE3-4    3.00 ± 0%
//    HashSHA3-4      6.00 ± 0%
// ----------------------------------------------------------------------------

func BenchmarkHashBLAKE3(b *testing.B) {
	// Ensure that hash algorithm is set to BLAKE3.
	oldHashAlgo := hasher.HashAlgo
	defer func() {
		hasher.HashAlgo = oldHashAlgo
	}()

	hasher.HashAlgo = hasher.HashAlgoBLAKE3

	const input = "This is a string"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)

		_, _ = hasher.Hash(r)
	}
}

func BenchmarkHashSHA3(b *testing.B) {
	// Ensure that hash algorithm is set to SHA3-512.
	oldHashAlgo := hasher.HashAlgo
	defer func() {
		hasher.HashAlgo = oldHashAlgo
	}()

	hasher.HashAlgo = hasher.HashAlgoSHA3_512

	const input = "This is a string"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)

		_, _ = hasher.Hash(r)
	}
}

package hasher_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/KEINOS/go-genrawid/pkg/hasher"
)

func ExampleCheckSum() {
	input := "1234567890"
	r := strings.NewReader(input)

	sumByte, err := hasher.CheckSum(r)
	if err != nil {
		log.Fatal(err)
	}

	// By default, the checksum algorithm is CRC32C with 4 bytes of length.
	fmt.Printf("Checksum algorithm: %v\n", hasher.ChkSumAlgo.String())
	fmt.Printf("Length of checksum: %v bytes\n", len(sumByte))
	fmt.Printf("Checksum value    : %x\n", sumByte)

	// Output:
	// Checksum algorithm: crc32
	// Length of checksum: 4 bytes
	// Checksum value    : f3dbd4fe
}

func ExampleCheckSum_xxhash() {
	// backup and defer restore checksum algorithm
	oldChkSumAlgo := hasher.ChkSumAlgo
	defer func() {
		hasher.ChkSumAlgo = oldChkSumAlgo
	}()

	// Change checksum algorithm from CRC32C to xxHash
	hasher.ChkSumAlgo = hasher.ChkSumXXHash

	input := "1234567890"
	r := strings.NewReader(input)

	sumByte, err := hasher.CheckSum(r)
	if err != nil {
		log.Fatal(err)
	}

	// By default, the checksum algorithm is CRC32C with 4 bytes of length.
	fmt.Printf("Checksum algorithm: %v\n", hasher.ChkSumAlgo.String())
	fmt.Printf("Length of checksum: %v bytes\n", len(sumByte))
	fmt.Printf("Checksum value    : %x\n", sumByte)

	// Output:
	// Checksum algorithm: xxhash
	// Length of checksum: 4 bytes
	// Checksum value    : a9d4d413
}

func ExampleHash() {
	input := "This is a string"
	r := strings.NewReader(input)

	// By default, the hash algorithm is BLAKE3 with 64 bytes of length.
	hashByte, err := hasher.Hash(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hash algorithm(BLAKE3): %v\n", hasher.HashAlgo.String())
	fmt.Printf("Length of full hash   : %v bytes\n", len(hashByte))
	fmt.Printf("First 32 bytes of hash: %x\n", hashByte[:32])

	// To change the hash algorithm set the hasher.HashAlgo global variable.
	// Here, we temporary set SHA3-512 as the hashing algorithm.
	hasher.HashAlgo = hasher.HashAlgoSHA3_512
	defer func() {
		// Restore to default after the test
		hasher.HashAlgo = hasher.HashAlgoBLAKE3
	}()

	// SHA3-512 returns 64 bytes of length hash as well but slower than BLAKE3
	hashByte, err = hasher.Hash(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hash algorithm(SHA3)  : %v\n", hasher.HashAlgo.String())
	fmt.Printf("Length of full hash   : %v bytes\n", len(hashByte))
	fmt.Printf("First 32 bytes of hash: %x\n", hashByte[:32])

	// Output:
	// Hash algorithm(BLAKE3): blake3
	// Length of full hash   : 64 bytes
	// First 32 bytes of hash: 718b749f12a61257438b2ea6643555fd995001c9d9ff84764f93f82610a780f2
	// Hash algorithm(SHA3)  : sha3-512
	// Length of full hash   : 64 bytes
	// First 32 bytes of hash: a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a6
}

func ExampleTChkSumAlgo_String() {
	for _, chkSumAlgo := range []hasher.TChkSumAlgo{
		hasher.ChkSumUnknown,
		hasher.ChkSumCRC32,
		hasher.ChkSumXXHash,
	} {
		fmt.Printf("Value: %#v, String: %s\n", chkSumAlgo, chkSumAlgo)
	}

	// Output:
	// Value: 0, String: unknown
	// Value: 1, String: crc32
	// Value: 2, String: xxhash
}

func ExampleTHashAlgo_String() {
	for _, chkSumAlgo := range []hasher.THashAlgo{
		hasher.HashAlgoUnknown,
		hasher.HashAlgoBLAKE3,
		hasher.HashAlgoSHA3_512,
	} {
		fmt.Printf("Value: %#v, String: %s\n", chkSumAlgo, chkSumAlgo)
	}

	// Output:
	// Value: 0, String: unknown
	// Value: 1, String: blake3
	// Value: 2, String: sha3-512
}

package hasher

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// The _sha3_512 returns the SHA3-512 hash of input. If lenOut is 0, the output
// length is 64 bytes.
func _sha3_512(input io.Reader, lenOut int) ([]byte, error) {
	lenMax := 64

	if lenOut == 0 {
		lenOut = hashLenDefault
	}

	if lenOut < 1 || lenOut > lenMax {
		return nil, errors.Errorf(
			"invalid output length. It must be between 1 and %d. Given length: %d",
			lenMax,
			lenOut,
		)
	}

	// Create new Hasher that has a digest size of 32 "bytes".
	sha3Hasher := sha3.New512()

	// Read the input data
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// sha3.state.Write panics if more data is written and never returns an error.
		_, _ = sha3Hasher.Write(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to read during scanning")
	}

	// Finalize the hash and return the digest.
	return sha3Hasher.Sum(nil)[:lenOut], nil
}

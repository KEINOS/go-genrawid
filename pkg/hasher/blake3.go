package hasher

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
	"github.com/zeebo/blake3"
)

// IoSeekStart is a copy of io.SeekStart to ease testing.
var IoSeekStart = io.SeekStart

// The _blake3 returns the BLAKE3 hash of input. If lenOut is 0, the output
// length is 64 bytes.
//
// It uses github.com/zeebo/blake3 package for BLAKE3 algorithm implementation
// of Go. This algorithm can generate a digest from 1 up to 8194 bytes of length.
func _blake3(input io.Reader, lenOut int) ([]byte, error) {
	lenMax := 8194

	if lenOut == 0 {
		lenOut = hashLenDefault
	}

	if lenOut > lenMax || lenOut < 1 {
		return nil, errors.Errorf(
			"invalid output length. It must be between 1 and 8194. Given length: %d\n",
			lenOut,
		)
	}

	// Create new Hasher that has a digest size of 32 "bytes".
	blake3Hasher := blake3.New()

	// Read the input data
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// The blake3.Hasher.Write never returns an error
		_, _ = blake3Hasher.Write(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to read during scanning")
	}

	// Finalize the hash and return the digest.
	// Digest takes a snapshot of the hash state and returns an object that can
	// be used to read and seek through 2^64 bytes of digest output.
	d := blake3Hasher.Digest()
	hashed := make([]byte, lenOut)

	if _, err := d.Seek(0, IoSeekStart); err != nil {
		return nil, errors.Wrap(err, "failed to set the position to seek")
	}

	// The blake3.Digest.Read always fills the entire buffer and never errors.
	_, _ = d.Read(hashed)

	return hashed, nil
}

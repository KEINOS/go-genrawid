package hasher

import (
	"io"

	"github.com/cespare/xxhash/v2"
	"github.com/pkg/errors"
)

// The _xxhash returns the 8 Byte/64 bit hash value of xxHash from the input with
// the length of lenOut.
//
// The lenOut must be between 0-8. If lenOut is 0(zero) then it will return with
// the max length 8 Byte.
func _xxhash(input io.Reader, lenOut int) ([]byte, error) {
	if err := preCheckXxHash(&lenOut); err != nil {
		return nil, errors.Wrap(err, "error on _xxhash.")
	}

	xxHasher := xxhash.New()
	buffer := make([]byte, 1) // set buffer size as 1 Byte to read

	for {
		n, err := input.Read(buffer) // read input to buffer
		if n > 0 {
			// xxhash.Digest.Write always returns len(b), nil and no error
			_, _ = xxHasher.Write(buffer)

			continue
		}

		if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "error during reading data")
		}

		if err == io.EOF || n == 0 {
			break
		}
	}

	return xxHasher.Sum(nil)[:lenOut], nil
}

func preCheckXxHash(lenOut *int) error {
	const maxLen = 8

	if *lenOut == 0 {
		*lenOut = maxLen
	}

	if *lenOut > maxLen || *lenOut < 0 {
		return errors.Errorf("lenOut is too long or short. must be between 0-8. given: %v\n", lenOut)
	}

	return nil
}

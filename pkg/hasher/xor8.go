package hasher

import (
	"io"

	"github.com/pkg/errors"
)

// The _xor8 returns a 8bit/1Byte XOR8 checksum from the input.
func _xor8(input io.Reader) ([]byte, error) {
	const mask = uint(0b11111111)

	sum := uint(0)
	buffer := make([]byte, 1)

	for {
		lenRead, err := input.Read(buffer)
		if lenRead > 0 {
			sum = (sum + uint(buffer[0])) & mask

			continue
		}

		if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "error during reading data")
		}

		if err == io.EOF || lenRead == 0 {
			break
		}
	}

	checksum := ((sum ^ mask) + 1) & mask

	return []byte{byte(checksum)}, nil
}

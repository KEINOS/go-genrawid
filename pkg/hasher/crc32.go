package hasher

import (
	"hash/crc32"
	"io"

	"github.com/pkg/errors"
)

func _crc32(input io.Reader) ([]byte, error) {
	crcTable := crc32.MakeTable(CRC32Poly)
	hash32 := crc32.New(crcTable)

	if _, err := io.Copy(hash32, input); err != nil {
		return nil, errors.Wrap(err, "failed to copy data to hasher")
	}

	return hash32.Sum(nil), nil
}

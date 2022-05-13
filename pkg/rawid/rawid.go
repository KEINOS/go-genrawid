package rawid

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// NewBase62 returns new rawid.ID object from base62 encoded string.
func NewBase62(base62 string) (ID, error) {
	i := new(big.Int)

	bInt, ok := i.SetString(base62, 62)
	if !ok {
		return nil, errors.Errorf("fail to decode Base62 input: %s", base62)
	}

	var result ID = bInt.Bytes()

	if len(result) > 8 {
		return nil, errors.Errorf(
			"over range. The given Base62 string has more than 8 bytes after decoding. input: %s",
			base62,
		)
	}

	return result, nil
}

// ----------------------------------------------------------------------------
//  Type: ID
// ----------------------------------------------------------------------------

// ID is the type of the returned value as a rawid. It contains methods to
// convert the rawid to other types.
type ID []byte

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// Base62 returns the rawid as a base62 encoded string.
// The characters used are "0-9, a-z, A-Z".
func (r ID) Base62() string {
	i := new(big.Int)

	i.SetBytes(r[:])

	return i.Text(62)
}

// Byte returns the rawid as a byte slice.
func (r ID) Byte() []byte {
	return r
}

// Dec returns the rawid as a signed decimal string.
// It is the most suitable type to use in SQLite3's rawid.
func (r ID) Dec() string {
	return fmt.Sprintf("%d", r.Int64())
}

// Hex returns the rawid as a hex string.
func (r ID) Hex() string {
	return fmt.Sprintf("%x", r)
}

// Int64 returns the rawid as a signed 64 bit integer.
func (r ID) Int64() int64 {
	return int64(r.UInt64())
}

// String returns the rawid as a string. It is useless for rawid.
func (r ID) String() string {
	return string(r)
}

// UDec returns the rawid as an unsigned decimal string (No plus or minus).
func (r ID) UDec() string {
	return fmt.Sprintf("%d", r.UInt64())
}

// UInt64 returns the rawid as an unsigned 64 bit integer.
func (r ID) UInt64() uint64 {
	return binary.BigEndian.Uint64(r)
}

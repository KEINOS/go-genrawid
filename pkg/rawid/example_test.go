package rawid_test

import (
	"fmt"
	"log"

	"github.com/KEINOS/go-genrawid/pkg/rawid"
)

// ----------------------------------------------------------------------------
//  Examples of Public Functions
// ----------------------------------------------------------------------------

func ExampleNewBase62() {
	input := "lYGhA16ahyf" // Max value in Base62

	rawID, err := rawid.NewBase62(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Base10:", rawID.Dec())
	fmt.Println("Base16:", rawID.Hex())
	fmt.Println("Base62:", rawID.Base62())

	// Output:
	// Base10: -1
	// Base16: ffffffffffffffff
	// Base62: lYGhA16ahyf
}

// ----------------------------------------------------------------------------
//  Examples of RawID methods
// ----------------------------------------------------------------------------

func ExampleID_Base62() {
	// the ID must be 8 byte length
	id := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	fmt.Println("Base10:", id.Dec())
	fmt.Println("Base62:", id.Base62())

	// Output:
	// Base10: -1
	// Base62: lYGhA16ahyf
}

// To get the rawid as a slice of bytes, use the ID.Bytes() method.
func ExampleID_Byte() {
	// 0x61 = 0d97, 0x62 = 0d98, ... ... , 0x68 = 0d104
	id := rawid.ID{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68}

	fmt.Println("Byte silice:", id.Byte())

	// Output:
	// Byte silice: [97 98 99 100 101 102 103 104]
}

func ExampleID_String() {
	id := rawid.ID{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68}

	// Type ID implements the Stringer interface. Note that both of the
	// below are equivalent.
	fmt.Println("Stringer:", id)
	fmt.Println("String  :", id.String())

	// Output:
	// Stringer: abcdefgh
	// String  : abcdefgh
}

// To get the rawid as 16 digit hexadecimal string, use the RawID.Hex() method.
//
// Note that the first zeros are not omitted and the string is always 16 digits.
func ExampleID_Hex() {
	idMin := rawid.ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMid := rawid.ID{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMax := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	fmt.Println("Hex string (min):", idMin.Hex())
	fmt.Println("Hex string (mid):", idMid.Hex())
	fmt.Println("Hex string (max):", idMax.Hex())

	// Output:
	// Hex string (min): 0000000000000000
	// Hex string (mid): 8000000000000000
	// Hex string (max): ffffffffffffffff
}

// To get the rawid as signed decimal string, use the RawID.Dec() method.
//
// Note that numbers are treated unsigned internally and the output is a
// signed number. Check the RawID.UDec method example as well.
func ExampleID_Dec() {
	idMin := rawid.ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMid := rawid.ID{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMax := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// Signed decimal
	fmt.Println("Dec string (min):", idMin.Dec())
	fmt.Println("Dec string (mid):", idMid.Dec())
	fmt.Println("Dec string (max):", idMax.Dec())

	// Output:
	// Dec string (min): 0
	// Dec string (mid): -9223372036854775808
	// Dec string (max): -1
}

// To get the rawid as an unsigned decimal string, use the ID.UDec() method.
// Check the ID.Dec method example as well.
func ExampleID_Dec_and_UDec() {
	idMin := rawid.ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMid := rawid.ID{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	idMax := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// Unsigned decimal
	fmt.Println("UDec string (min):", idMin.UDec())
	fmt.Println("UDec string (mid):", idMid.UDec())
	fmt.Println("UDec string (max):", idMax.UDec())

	// Output:
	// UDec string (min): 0
	// UDec string (mid): 9223372036854775808
	// UDec string (max): 18446744073709551615
}

// To get the rawid as an Int64 type, use the RawID.Int64() method.
// Check the RawID.UInt64 method example as well.
func ExampleID_Int64() {
	id := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// Signed int64
	fmt.Printf("Type: %T, Value: %v\n", id.Int64(), id.Int64())

	// Output:
	// Type: int64, Value: -1
}

// To get the rawid as an UInt64 type, use the RawID.UInt64() method.
// Check the RawID.Int64 method example as well.
func ExampleID_UInt64() {
	id := rawid.ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// Unsigned int64
	fmt.Printf("Type: %T, Value: %v\n", id.UInt64(), id.UInt64())

	// Output:
	// Type: uint64, Value: 18446744073709551615
}

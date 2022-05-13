package genrawid_test

import (
	"fmt"
	"log"
	"os"

	"github.com/KEINOS/go-genrawid"
)

// ----------------------------------------------------------------------------
//  Examples for public functions
// ----------------------------------------------------------------------------

func ExampleFromString() {
	input := "abcdefgh"

	rawid, err := genrawid.FromString(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rawid.Hex())
	fmt.Println(rawid.Dec())

	// Output:
	// ddaa2ac39b79058a
	// -2474118025671277174
}

func ExampleFromFile() {
	pathFile := "./testdata/msg.txt" // msg.txt ==> "abcdefgh"

	rawid, err := genrawid.FromFile(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rawid.Hex())
	fmt.Println(rawid.Dec())

	// Output:
	// ddaa2ac39b79058a
	// -2474118025671277174
}

// ExampleFromStdin
//
// Since we can not receive input from stdin during the example run, we mock
// the OsStdin, the copy of os.Stdin in the package, by assigning a dummy
// file as its input.
func ExampleFromStdin() {
	// Backup and defer recover the stdin
	oldStdin := genrawid.OsStdin
	defer func() {
		genrawid.OsStdin = oldStdin
	}()

	// Mock the stdin by assigning a file to stdin.
	pathFile := "./testdata/msg.txt" // msg.txt ==> "abcdefgh"

	osFile, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	genrawid.OsStdin = osFile // <-- mock!

	// Here is the actual example usage of FromStdin method to get the input from
	// stdin and generate a rawid.
	rawid, err := genrawid.FromStdin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rawid.Hex())
	fmt.Println(rawid.Dec())

	// Output:
	// ddaa2ac39b79058a
	// -2474118025671277174
}

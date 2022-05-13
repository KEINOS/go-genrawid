package main

import (
	"fmt"
	"os"

	"github.com/KEINOS/go-utiles/util"
	"github.com/spf13/pflag"
)

func usage() {
	fmt.Fprintln(os.Stderr, util.HereDoc(`
		genrawid - generates a unique consistent number from the imput.

		Usage:
		  genrawid [flags] [filepath | - ]

		  If "filepath" argument is "-" then it will read from the piped STDIN.
	`))

	fmt.Fprintln(os.Stderr, "Flags:")
	pflag.PrintDefaults()

	fmt.Fprintln(os.Stderr, util.HereDoc(`

		Example:
		  $ # To specify a file to generate its rawid.
		  $ genrawid /path/to/my/file.pdf

		  $ # Specify a file content via STDIN. The following two are equivalent.
		  $ genrawid - < /path/to/my/file.pdf
		  $ cat /path/to/my/file.pdf | genrawid -

		  $ # Specify text content via STDIN.
		  $ echo 'foo bar' | genrawid -

		  $ # Specify text content via argument.
		  $ genrawid -s "foo bar"

		  $ # Adds a line break to the output.
		  $ genrawid -s "foo bar" --new-line

		  $ # Verify if rawid is equivalent to the given rawid. It will exit with
		  $ # status 0 if matches, and 1 if not.
		  $ genrawid -s "foo bar" --verify "-7374369981397550869"

		About:
		  genrawid is a niche tool to generate a unique number from the input.

		  The returned value is a 64 bit/8 Byte signed decimals which is a
		  mixture of a hash and checksum values. The first 32 bit/4 Byte is the
		  hash and the second 32 bit/4 Byte is the checksum of the full hash
		  value. By default it uses BLAKE3-512 as a hash algorithm and CRC-32C
		  as a checksum.

		  The objective is to generate a 64 bit/8 Byte of rawid to use in SQLite3
		  as a fast KVS (Key-Value-Store) for CAS (Content-Addressable-Storage)
		  usage.

		  For more info see: https://github.com/KEINOS/go-genrawid

	`))
}

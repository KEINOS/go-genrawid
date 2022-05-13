package main

import (
	"fmt"
	"hash/crc32"

	"github.com/KEINOS/go-genrawid"
	"github.com/KEINOS/go-genrawid/pkg/hasher"
	"github.com/KEINOS/go-genrawid/pkg/rawid"
	"github.com/KEINOS/go-utiles/util"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	inStr    = "" // it holds the input string from the arg.
	inVerify = "" // it holds the gigen rawid to compare.
	lf       = "" // line-feed to use if set.
	pathFile = "" // file path to read if set.

	isBase62 bool // outputs the results in base62 if true.
	isFile   bool // read input from file.
	isHelp   bool // diplays help if true.
	isHex    bool // outputs the results in hex if true.
	isLF     bool // line breaks the output if true.
	isStdin  bool // receive input from STDIN if true.
	isString bool // receive input from command arg.
	isVerify bool // compares between the given rawid and calculated rawid.
)

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

func main() {
	util.ExitOnErr(Run())
}

// ----------------------------------------------------------------------------
//  Public Functions
// ----------------------------------------------------------------------------

// PreRun : parse flag and checks the option status.
func PreRun() error {
	// Set and parse flags
	setFlags()
	pflag.Parse()

	pflag.Usage = usage // set custom usage for help msg

	args := pflag.Args() // get non-flag command-line arguments.

	// Args length validation
	if len(args) > 1 {
		return errors.New("too many arguments. more than one file path given")
	}

	chkOptFile(args)  // file path check
	chkOptLineFeed()  // --new-line option check
	chkOptStdin(args) // - (stdin) option check
	chkOptString()    // --string option check
	chkOptVerify()    // --verify option check

	switch {
	case isHelp:
		return nil
	case isString:
		return nil
	case isStdin:
		return nil
	case isFile:
		return nil
	}

	return errors.New("error: missing arguments")
}

// Run is the actual function of the app.
func Run() (err error) {
	if err = PreRun(); err != nil {
		return errors.Wrap(err, "error during pre-run")
	}

	var id rawid.ID

	switch {
	case isHelp:
		pflag.Usage()

		return nil
	case isString:
		// FromString generates whater the input is
		id, _ = genrawid.FromString(inStr)
	case isFile:
		id, err = genrawid.FromFile(pathFile)
		if err != nil {
			return errors.Wrap(err, "failed to read from file")
		}
	case isStdin:
		id, err = genrawid.FromStdin()
		if err != nil {
			return errors.Wrap(err, "failed to read from STDIN")
		}
	}

	output := ""

	// Format output
	if isHex {
		output = fmt.Sprintf("0x%v%v", id.Hex(), lf)
	} else if isBase62 {
		output = fmt.Sprintf("%v%v", id.Base62(), lf)
	} else {
		output = fmt.Sprintf("%v%v", id.Dec(), lf)
	}

	if isVerify {
		if output != (inVerify + lf) {
			return errors.Errorf(
				"the two rawids did not match. Given: %v, Calculated: %v\n",
				inVerify,
				output,
			)
		}
	}

	// Print the calculated rawid
	fmt.Print(output)

	return nil
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

// Set flags to default values.
func setFlags() {
	// Default algorithm
	hasher.HashAlgo = hasher.HashAlgoBLAKE3
	hasher.ChkSumAlgo = hasher.ChkSumCRC32
	hasher.CRC32Poly = crc32.Castagnoli

	// Default flag values
	isString = false
	isStdin = false
	isHelp = false
	isHex = false
	isBase62 = false
	isLF = false
	inStr = ""
	inVerify = ""
	lf = ""

	if !pflag.Parsed() {
		pflag.BoolVar(&isBase62, "base62", false, "outputs the rawid in Base62 encoded string (uses: 0-9,a-z,A-Z)")
		pflag.BoolVarP(&isHelp, "help", "h", false, "displays this help")
		pflag.BoolVar(&isHex, "hex", false, "outputs the rawid in hex string")
		pflag.BoolVarP(&isLF, "new-line", "n", false, "line-feed/line-breaks after the output")
		pflag.StringVarP(&inStr, "string", "s", "", "provide the input via args")
		pflag.StringVar(&inVerify, "verify", "", "the rawid to verify")
	}
}

func chkOptFile(args []string) {
	if len(args) > 0 && args[0] != "-" {
		// Stdin and File input cannot coexist
		isStdin = false
		isFile = true
		pathFile = args[0]
	}
}

func chkOptLineFeed() {
	if isLF {
		// Use \n as a line break char
		lf = "\n"
	}
}

func chkOptStdin(args []string) {
	if len(args) > 0 && args[0] == "-" {
		// Stdin and File input cannot coexist
		isStdin = true
		isFile = false
	}
}

func chkOptString() {
	if inStr != "" {
		// Flag up if input was given via arg
		isString = true
	}
}

func chkOptVerify() {
	if inVerify != "" {
		// Flag up if rawid was provided to verify
		isVerify = true
	}
}

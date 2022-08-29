package main

import (
	"fmt"
	"hash/crc32"
	"os"

	"github.com/KEINOS/go-genrawid"
	"github.com/KEINOS/go-genrawid/pkg/hasher"
	"github.com/KEINOS/go-genrawid/pkg/rawid"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	inStr    string // it holds the input string from the arg.
	inVerify string // it holds the given rawid to compare.
	lineFeed string // line-feed to use if set.
	pathFile string // file path to read if set.

	isBase62 bool // outputs the results in base62 if true.
	isFast   bool // fast mode if true.
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
	ExitOnError(Run())
}

// ----------------------------------------------------------------------------
//  Public Functions
// ----------------------------------------------------------------------------

// OsExit is a copy of os.Exit() to ease testing.
var OsExit = os.Exit

// ExitOnError exits with status 1 if err is an error.
func ExitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		OsExit(1)
	}
}

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
	chkModeFast()     // --fast option check

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
//
//nolint:cyclop // complexity is 13 but it's okay.
func Run() error {
	err := PreRun()
	if err != nil {
		return errors.Wrap(err, "error during pre-run")
	}

	//nolint:varnamelen // allow short variable names for readability
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

	var output string

	// Format output
	switch {
	case isHex:
		output = fmt.Sprintf("0x%v%v", id.Hex(), lineFeed)
	case isBase62:
		output = fmt.Sprintf("%v%v", id.Base62(), lineFeed)
	default:
		output = fmt.Sprintf("%v%v", id.Dec(), lineFeed)
	}

	if isVerify {
		if output != (inVerify + lineFeed) {
			return errors.Errorf(
				"the two rawids did not match. Given: %v, Calculated: %v",
				inVerify,
				output,
			)
		}
	}

	// Print the calculated rawid
	//nolint:forbidigo // allow printing to stdout
	fmt.Print(output)

	return nil
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

func chkModeFast() {
	genrawid.IsModeFast = isFast
}

func chkOptFile(args []string) {
	// Stdin and File input cannot coexist
	if len(args) > 0 && args[0] != "-" {
		isStdin = false
		isFile = true
		pathFile = args[0]
	}
}

func chkOptLineFeed() {
	if isLF {
		// Use \n as a line break char
		lineFeed = "\n"
	}
}

func chkOptStdin(args []string) {
	// Stdin and File input cannot coexist
	if len(args) > 0 && args[0] == "-" {
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

// Set flag/option values to default.
func resetFlagValues() {
	isBase62 = false
	isFast = false
	isFile = false
	isHelp = false
	isHex = false
	isLF = false
	isStdin = false
	isString = false
	isVerify = false

	inStr = ""
	inVerify = ""
	lineFeed = ""
	pathFile = ""
}

// Set flags to default values.
func setFlags() {
	// Set default algorithm
	hasher.HashAlgo = hasher.HashAlgoBLAKE3
	hasher.ChkSumAlgo = hasher.ChkSumCRC32
	hasher.CRC32Poly = crc32.Castagnoli

	// Set default mode
	genrawid.IsModeFast = false

	// Set default flag values
	resetFlagValues()

	// Initialize flags before parsing
	if !pflag.Parsed() {
		pflag.BoolVar(&isBase62, "base62", false, "outputs the rawid in Base62 encoded string (uses: 0-9,a-z,A-Z)")
		pflag.BoolVarP(&isHelp, "help", "h", false, "displays this help")
		pflag.BoolVarP(&isFast, "fast", "f", false, "fast mode (uses: XOR16 for checksum)")
		pflag.BoolVar(&isHex, "hex", false, "outputs the rawid in hex string")
		pflag.BoolVarP(&isLF, "new-line", "n", false, "line-feed/line-breaks after the output")
		pflag.StringVarP(&inStr, "string", "s", "", "provide the input via args")
		pflag.StringVar(&inVerify, "verify", "", "the rawid to verify")
	}

	// Hide 'fast' flag.
	// It wasn't fast enough as expected. So, currently it's disabled.
	_ = pflag.CommandLine.MarkHidden("fast")
	_ = pflag.CommandLine.MarkHidden("f")
}

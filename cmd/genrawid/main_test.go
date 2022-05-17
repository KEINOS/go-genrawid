package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-genrawid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

// ============================================================================
//  Tests
// ============================================================================

// ----------------------------------------------------------------------------
//  Golden Cases
// ----------------------------------------------------------------------------

func Test_main_golden_base62(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{
		"--base62",
		"../../testdata/msg.txt",
	})
	defer deferRecover()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "j1UNoJA6ku6"
	actual := out
	assert.Equal(t, expect, actual)
}

func Test_main_golden_file(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{
		"../../testdata/msg.txt",
	})
	defer deferRecover()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "-2474118025671277174"
	actual := out
	assert.Equal(t, expect, actual)
}

func Test_main_golden_help(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{"-h"})
	defer deferRecover()

	out := capturer.CaptureStderr(func() {
		main()
	})

	assert.Contains(t, out, "Usage:")
	assert.Contains(t, out, "Flags:")
}

func Test_main_golden_hex_and_newline(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{
		"--new-line",
		"--hex",
		"../../testdata/msg.txt",
	})
	defer deferRecover()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "0xddaa2ac39b79058a\n"
	actual := out
	assert.Equal(t, expect, actual)
}

func Test_main_golden_stdin(t *testing.T) {
	// Mock stdin
	deferRecover := mockSTDIN(t, "abcdefgh")
	defer deferRecover()

	// Set args
	recoverArgs := setDummyArgs(t, []string{"-"})
	defer recoverArgs()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "-2474118025671277174"
	actual := out
	assert.Equal(t, expect, actual)
}

func Test_main_golden_string(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{"-s", "abcdefgh"})
	defer deferRecover()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "-2474118025671277174"
	actual := out
	assert.Equal(t, expect, actual)
}

func Test_main_golden_verify(t *testing.T) {
	// Set args
	deferRecover := setDummyArgs(t, []string{
		"-s",
		"abcdefgh",
		"--verify",
		"-2474118025671277174",
	})
	defer deferRecover()

	out := capturer.CaptureStdout(func() {
		main()
	})

	expect := "-2474118025671277174"
	actual := out
	assert.Equal(t, expect, actual)
}

// ----------------------------------------------------------------------------
//  Error Cases
// ----------------------------------------------------------------------------

func Test_main_missing_args(t *testing.T) {
	// Set empty args and defer recover
	recoverArgs := setDummyArgs(t, []string{})
	defer recoverArgs()

	// Mock os.Exit to capture exit status
	var status int

	recoverOsExit := captureExitStatus(t, &status)
	defer recoverOsExit()

	// Capture error
	out := capturer.CaptureOutput(func() {
		main()
	})

	assert.Equal(t, 1, status, "it should exit with status 1 on error")
	assert.Contains(t, out, "error: missing arguments")
}

func Test_main_path_was_dir(t *testing.T) {
	// Set empty args and defer recover
	recoverArgs := setDummyArgs(t, []string{
		t.TempDir(),
	})
	defer recoverArgs()

	// Mock os.Exit to capture exit status
	var status int

	recoverOsExit := captureExitStatus(t, &status)
	defer recoverOsExit()

	// Capture error
	out := capturer.CaptureStderr(func() {
		main()
	})

	assert.Equal(t, 1, status, "it should exit with status 1 on error")
	assert.Contains(t, out, "failed to read from file")
	assert.Contains(t, out, "failed to read during scanning")
}

func Test_main_stdin_error(t *testing.T) {
	// Backup the stdin
	oldOsStdin := genrawid.OsStdin
	defer func() {
		genrawid.OsStdin = oldOsStdin
	}()

	p, err := os.Open(t.TempDir())
	require.NoError(t, err)

	genrawid.OsStdin = p

	// Set empty args and defer recover
	recoverArgs := setDummyArgs(t, []string{"-"})
	defer recoverArgs()

	// Mock os.Exit to capture exit status
	var status int

	recoverOsExit := captureExitStatus(t, &status)
	defer recoverOsExit()

	// Capture error
	out := capturer.CaptureStderr(func() {
		main()
	})

	assert.Equal(t, 1, status, "it should exit with status 1 on error")

	assert.Contains(t, out, "failed to read from STDIN")
	assert.Contains(t, out, "failed to read during scanning")
}

func Test_main_too_many_args(t *testing.T) {
	// Set empty args and defer recover
	recoverArgs := setDummyArgs(t, []string{
		"../../testdata/msg.txt",
		"../../testdata/msg.txt",
	})
	defer recoverArgs()

	// Mock os.Exit to capture exit status
	var status int

	recoverOsExit := captureExitStatus(t, &status)
	defer recoverOsExit()

	// Capture error
	out := capturer.CaptureStderr(func() {
		main()
	})

	assert.Equal(t, 1, status, "it should exit with status 1 on error")

	contains := "too many arguments. more than one file path given"
	assert.Contains(t, out, contains)
}

func Test_main_verify_error(t *testing.T) {
	// Set empty args and defer recover
	recoverArgs := setDummyArgs(t, []string{
		"-s",
		"abcdefgh",
		"--verify",
		"foovar", // invalid rawid to verify
	})
	defer recoverArgs()

	// Mock os.Exit to capture exit status
	var status int

	recoverOsExit := captureExitStatus(t, &status)
	defer recoverOsExit()

	// Capture error
	out := capturer.CaptureStderr(func() {
		main()
	})

	assert.Equal(t, 1, status, "it should exit with status 1 on error")
	assert.Contains(t, out, "the two rawids did not match")
}

// ============================================================================
//  Helper Functions
// ============================================================================

func captureExitStatus(t *testing.T, status *int) func() {
	t.Helper()

	oldOsExit := OsExit

	OsExit = func(code int) {
		fmt.Fprintln(os.Stderr, "exit status:", code)
		*status = code
	}

	return func() {
		OsExit = oldOsExit
	}
}

func mockSTDIN(t *testing.T, input string) func() {
	t.Helper()

	// Temprorary save input to file
	pathFile := filepath.Join(t.TempDir(), t.Name())

	err := os.WriteFile(pathFile, []byte(input), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	// Backup the stdin
	oldOsStdin := genrawid.OsStdin

	// Mock the stdin by assigning a temp file to stdin.
	osFile, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	genrawid.OsStdin = osFile // <-- mock!

	return func() {
		osFile.Close()

		genrawid.OsStdin = oldOsStdin // recover the stdin
	}
}

func setDummyArgs(t *testing.T, args []string) func() {
	t.Helper()

	// Backup args
	oldOsArgs := os.Args

	// Mock args
	os.Args = []string{t.Name()}
	if len(args) > 0 {
		os.Args = append(os.Args, args...)
	}

	return func() {
		os.Args = oldOsArgs // restore args
	}
}

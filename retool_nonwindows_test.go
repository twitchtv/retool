// +build !windows

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Go builds files on windows with an '.exe' suffix. Everywhere else, there's no
// suffix.
const osBinSuffix = ""

// Test that we correctly preserve .c and .h files by running a test against a
// command that uses go-sqlite3.
//
// This test can only be run on non-windows platforms because go-sqlite3 cannot
// be built on windows.
func TestCSourceFilePreservation(t *testing.T) {
	t.Parallel()

	retool, err := buildRetool()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(filepath.Dir(retool))
	}()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to make temp dir: %s", err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	runRetoolCmd(t, dir, retool, "add", "github.com/spenczar/sqlite_retool_testcmd", "origin/master")
	runRetoolCmd(t, dir, retool, "build")
}

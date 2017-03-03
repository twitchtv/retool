package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

// build builds retool in a temporary directory and returns the path to the built binary
func build() (string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", errors.Wrap(err, "unable to create temporary build directory")
	}
	output := filepath.Join(dir, "retool")
	cmd := exec.Command("go", "build", "-o", output, ".")
	_, err = cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "unable to build retool binary")
	}
	return output, nil
}

func TestRetoolCachePollution(t *testing.T) {
	retool, err := build()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(retool)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to make temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	// This should fail because this version of mockery has an import line that points to uber's
	// internal repo, which can't be reached:
	cmd := exec.Command(retool, "-base-dir", dir, "add",
		"github.com/vektra/mockery/cmd/mockery", "d895b9fcc32730719faaccd7840ad7277c94c2d0",
	)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err == nil {
		t.Fatal("expected error when adding mockery at broken commit d895b9, but got no error")
	}

	// Now, without cleaning the cache, try again on a healthy commit. In
	// ff9a1fda7478ede6250ee3c7e4ce32dc30096236 of retool and earlier, this would still fail because
	// the cache would be polluted with a bad source tree.
	cmd = exec.Command(retool, "-base-dir", dir, "add",
		"github.com/vektra/mockery/cmd/mockery", "origin/master",
	)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		t.Fatalf("expected no error when adding mockery at broken commit d895b9, but got this:\n%s", string(err.(*exec.ExitError).Stderr))
	}
}

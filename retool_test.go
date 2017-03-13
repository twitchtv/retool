package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

// buildRetool builds retool in a temporary directory and returns the path to the built binary
func buildRetool() (string, error) {
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

func TestRetool(t *testing.T) {
	// These integration tests require more than most go tests: they require a go compiler to build
	// retool, a working version of git to perform retool's operations, and network access to do the
	// git fetches.
	retool, err := buildRetool()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(retool)
	}()

	t.Run("cache pollution", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("unable to make temp dir: %s", err)
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

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
	})

	t.Run("version", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("unable to make temp dir: %s", err)
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		// Should work even in a directory without tools.json
		cmd := exec.Command(retool, "version")
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("expected no errors when using retool version, have this:\n%s", string(out))
		}
		if want := fmt.Sprintf("retool %s", version); string(out) != want {
			t.Errorf("have=%q, want=%q", string(out), want)
		}
	})

	t.Run("build", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("unable to make temp dir: %s", err)
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		cmd := exec.Command(retool, "-base-dir", dir, "add",
			"github.com/twitchtv/retool", "origin/master",
		)
		cmd.Dir = dir
		_, err = cmd.Output()
		if err != nil {
			t.Fatalf("expected no errors when using retool add, have this:\n%s", string(err.(*exec.ExitError).Stderr))
		}

		// Suppose we only have _tools/src available. Does `retool build` work?
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "bin"))
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "pkg"))
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "manifest.json"))

		cmd = exec.Command(retool, "-base-dir", dir, "build")
		cmd.Dir = dir
		_, err = cmd.Output()
		if err != nil {
			t.Fatalf("expected no errors when using retool build, have this:\n%s", string(err.(*exec.ExitError).Stderr))
		}

		// Now the binary should be installed
		_, err = os.Stat(filepath.Join(dir, "_tools", "bin", "retool"))
		if err != nil {
			t.Fatalf("unable to stat _tools/bin/retool after calling retool build: %s", err)
		}
	})

	t.Run("dep_added", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("unable to make temp dir: %s", err)
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		// Use a package which used to have a dependency (in this case, one on
		// github.com/spenczar/retool_test_lib), but doesn't have that dependency for HEAD of
		// origin/master today.
		cmd := exec.Command(retool, "-base-dir", dir, "add",
			"github.com/spenczar/retool_test_app", "origin/has_dep",
		)
		cmd.Dir = dir
		_, err = cmd.Output()
		if err != nil {
			t.Fatalf("expected no errors when using retool add, have this:\n%s", string(err.(*exec.ExitError).Stderr))
		}
	})
}

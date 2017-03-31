package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRetool(t *testing.T) {
	// These integration tests require more than most go tests: they require a go compiler to build
	// retool, a working version of git to perform retool's operations, and network access to do the
	// git fetches.
	retool, err := buildRetool()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(filepath.Dir(retool))
	}()

	t.Run("cache pollution", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()

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
		runRetoolCmd(t, dir, retool, "add", "github.com/vektra/mockery/cmd/mockery", "origin/master")
	})

	t.Run("version", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()

		// Should work even in a directory without tools.json
		out := runRetoolCmd(t, dir, retool, "version")
		if want := fmt.Sprintf("retool %s", version); string(out) != want {
			t.Errorf("have=%q, want=%q", string(out), want)
		}
	})

	t.Run("build", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()
		runRetoolCmd(t, dir, retool, "add", "github.com/twitchtv/retool", "origin/master")

		// Suppose we only have _tools/src available. Does `retool build` work?
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "bin"))
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "pkg"))
		_ = os.RemoveAll(filepath.Join(dir, "_tools", "manifest.json"))

		runRetoolCmd(t, dir, retool, "build")

		// Now the binary should be installed
		assertBinInstalled(t, dir, "retool")

		// Legal files should be kept around
		_, err = os.Stat(filepath.Join(dir, "_tools", "retool", "src", "github.com", "twitchtv", "retool", "LICENSE"))
		if err != nil {
			t.Error("missing license file")
		}
	})

	t.Run("dep_added", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()

		// Use a package which used to have a dependency (in this case, one on
		// github.com/spenczar/retool_test_lib), but doesn't have that dependency for HEAD of
		// origin/master today.
		runRetoolCmd(t, dir, retool, "add", "github.com/spenczar/retool_test_app", "origin/has_dep")
	})

	t.Run("multiple_deps_added", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()

		// Use a package which used to have a dependency (in this case, one on
		// github.com/spenczar/retool_test_lib), but doesn't have that dependency for HEAD of
		// origin/master today.
		runRetoolCmd(t, dir, retool, "add", "github.com/spenczar/retool_test_app", "origin/has_dep")
		runRetoolCmd(t, dir, retool, "add", "github.com/twitchtv/retool", "origin/master")

		// Now the binary should be installed
		assertBinInstalled(t, dir, "retool")
		assertBinInstalled(t, dir, "retool_test_app")

	})

	t.Run("clean", func(t *testing.T) {
		// Clean should be a noop, but kept around for compatibility
		cmd := exec.Command(retool, "clean")
		_, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				t.Fatalf("expected no errors when using retool clean, have this:\n%s", string(exitErr.Stderr))
			} else {
				t.Fatalf("unexpected err when running %q: %q", strings.Join(cmd.Args, " "), err)
			}
		}
	})

	t.Run("do", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()

		runRetoolCmd(t, dir, retool, "add", "github.com/twitchtv/retool", "v1.0.1")
		output := runRetoolCmd(t, dir, retool, "do", "retool", "version")
		if want := "retool v1.0.1"; output != want {
			t.Errorf("have=%q, want=%q", output, want)
		}
	})

	t.Run("upgrade", func(t *testing.T) {
		dir, cleanup := setupTempDir(t)
		defer cleanup()
		runRetoolCmd(t, dir, retool, "add", "github.com/twitchtv/retool", "v1.0.1")
		runRetoolCmd(t, dir, retool, "upgrade", "github.com/twitchtv/retool", "v1.0.3")
		out := runRetoolCmd(t, dir, retool, "do", "retool", "version")
		if want := "retool v1.0.3"; string(out) != want {
			t.Errorf("have=%q, want=%q", string(out), want)
		}
	})
	t.Run("gometalinter exemption", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("unable to make temp dir: %s", err)
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		runRetoolCmd(t, dir, retool, "add", "github.com/alecthomas/gometalinter", "origin/master")
		runRetoolCmd(t, dir, retool, "do", "gometalinter", "--install")

		// Create a dummy go file so gometalinter runs. If we don't do this,
		// gometalinter will exit without doing any work, and we'll get a false
		// positive.
		//
		// The file will be removed with the deferred os.RemoveAll(dir) call, no
		// need to remove it here.
		f, err := os.Create(filepath.Join(dir, "main.go"))
		if err != nil {
			t.Fatalf("unable to create file for gometalinter to run against: %s", err)
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				t.Errorf("unable to close gometalinter test file: %s", closeErr)
			}
		}()
		_, err = io.WriteString(f, `package main

func main() {}`)
		if err != nil {
			t.Fatalf("unable to write gometalinter test file: %s", err)
		}

		// If gometalinter can't find its tools, it will exit with code 2.
		runRetoolCmd(t, dir, retool, "do", "gometalinter", ".")

		// Make sure gometalinter installs to the tool directory, not to the global
		// GOPATH.
		assertBinInstalled(t, dir, "structcheck")
	})
}

func runRetoolCmd(t *testing.T, dir, retool string, args ...string) (output string) {
	args = append([]string{"-base-dir", dir}, args...)
	cmd := exec.Command(retool, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Fatalf("command %q failed, stderr:\n%s\n\nstdout:%s", "retool "+strings.Join(cmd.Args[1:], " "), string(exitErr.Stderr), string(out))
		} else {
			t.Fatalf("unexpected err when running %q: %q", strings.Join(cmd.Args, " "), err)
		}
	}
	return string(out)
}

func setupTempDir(t *testing.T) (dir string, cleanup func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to make temp dir: %s", err)
	}

	cleanup = func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("unable to clean up temp dir: %s", err)
		}
	}

	return dir, cleanup
}

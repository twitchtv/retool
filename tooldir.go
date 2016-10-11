package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	toolDirName = "_tools"
)

var (
	baseDir = flag.String("base-dir", "",
		"Path of project root.  If not specified, the working directory is used.")
	toolDir = flag.String("tool-dir", "",
		"Path where tools are stored.  The default value is the subdirectory of -base-dir named '_tools'.")

	// These globals are set by ensureTooldir() after factoring in the flags above.
	baseDirPath string
	toolDirPath string
)

func ensureTooldir() error {
	var err error

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get working directory")
	}

	baseDirPath = *baseDir
	if baseDirPath == "" {
		baseDirPath = cwd
	}

	toolDirPath = *toolDir
	if toolDirPath == "" {
		toolDirPath = filepath.Join(baseDirPath, toolDirName)
	}

	verbosef("base dir: %v\n", baseDirPath)
	verbosef("tool dir: %v\n", toolDirPath)

	stat, err := os.Stat(toolDirPath)
	switch {
	case os.IsNotExist(err):
		err = os.Mkdir(toolDirPath, 0777)
		if err != nil {
			return errors.Wrap(err, "unable to create tooldir")
		}
	case err != nil:
		return errors.Wrap(err, "unable to stat tool directory")
	case !stat.IsDir():
		return errors.New("tool directory already exists, but it is not a directory; you can use -tool-dir to change where tools are saved")
	}

	err = ioutil.WriteFile(path.Join(toolDirPath, ".gitignore"), gitignore, 0664)
	if err != nil {
		errors.Wrap(err, "unable to update .gitignore")
	}

	return nil
}

var gitignore = []byte(strings.TrimLeft(`
bin/
pkg/
manifest.json
`, "\n"))

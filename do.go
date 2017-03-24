package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func setPath() error {
	prevpath := os.Getenv("PATH")
	newPath := filepath.Join(toolDirPath, "bin") + string(os.PathListSeparator) + prevpath
	return os.Setenv("PATH", newPath)
}

func setGoEnv() error {
	newGoBin := filepath.Join(toolDirPath, "bin")
	if err := os.Setenv("GOBIN", newGoBin); err != nil {
		return err
	}

	prevGoPath := os.Getenv("GOPATH")
	newGoPath := toolDirPath + string(os.PathListSeparator) + prevGoPath
	return os.Setenv("GOPATH", newGoPath)
}

func do() {
	args := positionalArgs
	if len(args) == 0 {
		fatal("no command passed to retool do", nil)
	}

	if err := setPath(); err != nil {
		fatal("unable to set PATH", err)
	}
	if err := setGoEnv(); err != nil {
		fatal("unable to set up go environment variables", err)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		msg := "failed on '" + strings.Join(args, " ") + "'"
		fatal(msg, err)
	}
}

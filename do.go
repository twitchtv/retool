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

func setGoBin() error {
	newGoBin := filepath.Join(toolDirPath, "bin")
	return os.Setenv("GOBIN", newGoBin)
}

func do() {
	args := positionalArgs
	if len(args) == 0 {
		fatal("no command passed to retool do", nil)
	}

	if err := setPath(); err != nil {
		fatal("unable to set PATH", err)
	}
	if err := setGoBin(); err != nil {
		fatal("unable to set GOBIN", err)
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

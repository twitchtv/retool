package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func setPath() (unset func()) {
	prevpath := os.Getenv("PATH")
	newPath := filepath.Join(toolDirPath, "bin") + string(os.PathListSeparator) + prevpath
	_ = os.Setenv("PATH", newPath)
	return func() {
		_ = os.Setenv("PATH", prevpath)
	}
}

func setGoEnv() (unset func()) {
	prevGoPath, goPathWasSet := os.LookupEnv("GOPATH")
	var newGoPath string
	if goPathWasSet {
		newGoPath = prevGoPath + string(os.PathListSeparator) + toolDirPath
	} else {
		newGoPath = toolDirPath
	}
	_ = os.Setenv("GOPATH", newGoPath)

	prevGoBin, goBinWasSet := os.LookupEnv("GOBIN")
	newGoBin := filepath.Join(toolDirPath, "bin")
	_ = os.Setenv("GOBIN", newGoBin)

	return func() {
		if goPathWasSet {
			_ = os.Setenv("GOPATH", prevGoPath)
		} else {
			_ = os.Unsetenv("GOPATH")
		}
		if goBinWasSet {
			_ = os.Setenv("GOBIN", prevGoBin)
		} else {
			_ = os.Unsetenv("GOBIN")
		}
	}
}

func do() {
	args := positionalArgs
	if len(args) == 0 {
		fatal("no command passed to retool do", nil)
	}

	unsetPath := setPath()
	defer unsetPath()

	unsetGoEnv := setGoEnv()
	defer unsetGoEnv()

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

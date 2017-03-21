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

func do() {
	args := positionalArgs
	if len(args) == 0 {
		fatal("no command passed to retool do", nil)
	}

	defer setPath()()

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

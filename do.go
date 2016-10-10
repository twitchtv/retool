package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

// XXX: This does not actually set the return value ('unset').
func setPath() (unset func()) {
	prevpath := os.Getenv("PATH")
	newPath := path.Join(toolDirPath, "bin") + ":" + prevpath
	os.Setenv("PATH", newPath)
	return func() {
		os.Setenv("PATH", prevpath)
	}
}

func do(t *tool) {
	args := t.PositionalArgs
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

package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func setPath(cmd *exec.Cmd) {
	prevpath := os.Getenv("PATH")
	newPath := path.Join(tooldir, "bin") + ":" + prevpath
	setEnvVar(cmd, "PATH", newPath)
}

func trimArgs() []string {
	for i := 0; i < len(os.Args)-2; i++ {
		if os.Args[i] == "retool" && os.Args[i+1] == "do" {
			return os.Args[i+2:]
		}
	}
	return nil
}

func do() {
	args := trimArgs()
	if args == nil {
		fatal("no command passed to retool do", nil)
	}

	cmd := exec.Command(args[0], args[1:]...)

	setPath(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		msg := "failed on '" + strings.Join(args, " ") + "'"
		fatal(msg, err)
	}
}

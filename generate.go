package main

import (
	"os"
	"os/exec"
	"path"
)

func setPath(cmd *exec.Cmd) {
	prevpath := os.Getenv("PATH")
	newPath := path.Join(tooldir, "bin") + ":" + prevpath
	setEnvVar(cmd, "PATH", newPath)
}

func generate() {
	cmd := exec.Command("go", "generate", "./...")
	setGopath(cmd)
	setEnvVar(cmd, "PATH", path.Join(tooldir, "bin"))
	_, err := cmd.Output()
	if err != nil {
		fatal("failed on go generate", err)
	}
}

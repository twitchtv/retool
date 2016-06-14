package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type tool struct {
	Repository string // eg "github.com/tools/godep"
	Commit     string // eg "3020345802e4bff23902cfc1d19e90a79fae714e"
	ref        string // eg "origin/master"
}

func (t *tool) path() string {
	return path.Join(tooldir, "src", t.Repository)
}

func (t *tool) executable() string {
	return path.Base(t.Repository)
}

func setEnvVar(cmd *exec.Cmd, key, val string) {
	var env []string
	if cmd.Env != nil {
		env = cmd.Env
	} else {
		env = os.Environ()
	}

	envSet := false
	for i, envVar := range env {
		if strings.HasPrefix(envVar, key+"=") {
			env[i] = key + "=" + val
			envSet = true
		}
	}
	if !envSet {
		env = append(cmd.Env, key+"="+val)
	}

	cmd.Env = env
}

func setGopath(cmd *exec.Cmd) {
	setEnvVar(cmd, "GOPATH", tooldir)
}

func get(t *tool) error {
	log("downloading " + t.Repository)
	cmd := exec.Command("go", "get", "-d", t.Repository)
	setGopath(cmd)
	_, err := cmd.Output()
	return err
}

func setVersion(t *tool) error {
	log("setting version for " + t.Repository)
	cmd := exec.Command("git", "fetch")
	cmd.Dir = t.path()
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// If we have a symbolic reference, parse it
	if t.ref != "" {
		log(fmt.Sprintf("parsing revision %q", t.ref))
		cmd = exec.Command("git", "rev-parse", t.ref)
		cmd.Dir = t.path()
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		t.Commit = strings.TrimSpace(string(out))
		log(fmt.Sprintf("parsed as %q", t.Commit))
	}

	cmd = exec.Command("git", "checkout", t.Commit)
	cmd.Dir = t.path()
	_, err = cmd.Output()
	return err
}

func installBin(t *tool) error {
	log("installing " + t.Repository)
	cmd := exec.Command("go", "install", t.Repository)
	setGopath(cmd)
	_, err := cmd.Output()
	return err
}

func install(t *tool) error {
	err := get(t)
	if err != nil {
		fatalExec("go get -d "+t.Repository, err)
	}

	err = setVersion(t)
	if err != nil {
		fatalExec("git checkout "+t.Commit, err)
	}

	err = installBin(t)
	if err != nil {
		fatalExec("go install "+t.Repository, err)
	}

	return nil
}

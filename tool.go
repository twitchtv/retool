package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type tool struct {
	Repository string // eg "github.com/tools/godep"
	Commit     string // eg "3020345802e4bff23902cfc1d19e90a79fae714e"
	ref        string // eg "origin/master"
	Fork       string `json:"Fork,omitempty"` // eg "code.jusin.tv/twitch/godep"
}

func (t *tool) path() string {
	return filepath.Join(t.gopath(), "src", t.Repository)
}

func (t *tool) gopath() string {
	s := strings.Split(t.Repository, "/")
	return filepath.Join(toolDirPath, s[len(s)-1])
}

func (t *tool) executable() string {
	return filepath.Base(t.Repository)
}

func setEnvVar(cmd *exec.Cmd, key, val string) {
	var env []string
	// test if we added a var to this command before
	if cmd.Env != nil {
		env = cmd.Env
	} else {
		env = os.Environ()
	}

	// if we already have a key set for this var mutate it
	envSet := false
	for i, envVar := range env {
		if strings.HasPrefix(envVar, key+"=") {
			env[i] = key + "=" + val
			envSet = true
		}
	}
	// otherwise add it to our list
	if !envSet {
		env = append(env, key+"="+val)
	}

	// set the env
	cmd.Env = env
}

func get(t *tool) error {
	log("downloading " + t.Repository)
	cmd := exec.Command("go", "get", "-d", t.Repository)
	setEnvVar(cmd, "GOBIN", filepath.Join(toolDirPath, "bin"))
	setEnvVar(cmd, "GOPATH", t.gopath())
	_, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to 'go get' tool")
	}
	return err
}

func setVersion(t *tool) error {
	// If we're using a fork, add it
	if t.Fork != "" {
		cmd := exec.Command("git", "remote", "rm", "fork")
		cmd.Dir = t.path()
		_, err := cmd.Output()
		if err != nil {
			return err
		}

		cmd = exec.Command("git", "remote", "add", "-f", "fork", t.Fork)
		cmd.Dir = t.path()
		_, err = cmd.Output()
		if err != nil {
			return err
		}
	}

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

		var out []byte
		out, err = cmd.Output()
		if err != nil {
			return err
		}
		t.Commit = strings.TrimSpace(string(out))
		log(fmt.Sprintf("parsed as %q", t.Commit))
	}

	cmd = exec.Command("git", "checkout", t.Commit)
	cmd.Dir = t.path()
	_, err = cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to 'git checkout' tool")
	}

	// Re-run 'go get' in case the new version has a different set of dependencies.
	cmd = exec.Command("go", "get", "-d", t.Repository)
	setEnvVar(cmd, "GOBIN", filepath.Join(toolDirPath, "bin"))
	setEnvVar(cmd, "GOPATH", t.gopath())
	_, err = cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to 'go get' tool")
	}
	return err
}

func download(t *tool) error {
	err := get(t)
	if err != nil {
		fatalExec("go get -d "+t.Repository, err)
	}

	err = setVersion(t)
	if err != nil {
		fatalExec("git checkout "+t.Commit, err)
	}

	return nil
}

func install(t *tool) error {
	log("installing " + t.Repository)

	err := ioutil.WriteFile(path.Join(t.gopath(), ".gitignore"), gitignore, 0664)
	if err != nil {
		return errors.Wrap(err, "unable to update .gitignore")
	}

	cmd := exec.Command("go", "install", t.Repository)
	setEnvVar(cmd, "GOBIN", filepath.Join(toolDirPath, "bin"))
	setEnvVar(cmd, "GOPATH", t.gopath())
	_, err = cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to 'go install' tool")
	}
	return err
}

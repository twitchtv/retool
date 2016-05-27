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
}

var tools = []tool{
	// tool{
	// 	Executable: "counterfeiter",
	// 	Repository: "github.com/maxbrunsfeld/counterfeiter",
	// 	Commit:     "733ce8507ced45a0e8ae8bb0331b538a69cb1433",
	// },
	tool{
		Repository: "github.com/tools/godep",
		Commit:     "3020345802e4bff23902cfc1d19e90a79fae714e",
	},
}

func (t tool) path() string {
	return path.Join(tooldir, "src", t.Repository)
}

func (t tool) executable() string {
	return path.Base(t.Repository)
}

func (t tool) downloaded() (bool, error) {
	_, err := os.Stat(path.Join(tooldir, "bin", t.executable()))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t tool) correctVersion() (bool, error) {
	have, err := getVersion(t)
	if err != nil {
		return false, err
	}
	return have == t.Commit, nil
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

func get(t tool) error {
	log("downloading " + t.Repository)
	cmd := exec.Command("go", "get", "-d", t.Repository)
	setGopath(cmd)
	_, err := cmd.Output()
	return err
}

func setVersion(t tool) error {
	log("setting version for " + t.Repository)
	cmd := exec.Command("git", "fetch")
	cmd.Dir = t.path()
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "checkout", t.Commit)
	cmd.Dir = t.path()
	_, err = cmd.Output()
	return err
}

func getVersion(t tool) (string, error) {
	cmd := exec.Command("git", "rev-parse", "head")
	cmd.Dir = t.path()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func installBin(t tool) error {
	log("installing " + t.Repository)
	cmd := exec.Command("go", "install", t.Repository)
	setGopath(cmd)
	_, err := cmd.Output()
	return err
}

func install(t tool) error {
	downloaded, err := t.downloaded()
	if err != nil {
		fatal(fmt.Sprintf("unable to check if %s is downloaded", t.Repository), err)
	}
	if !downloaded {
		err := get(t)
		if err != nil {
			fatalExec("go get -d "+t.Repository, err)
		}
	}

	correctVersion, err := t.correctVersion()
	if err != nil {
		fatal(fmt.Sprintf("unable to check version of %s", t.Repository), err)
	}
	if !correctVersion {
		err = setVersion(t)
		if err != nil {
			fatalExec("git checkout "+t.Commit, err)
		}
	}

	if !downloaded || !correctVersion {
		err = installBin(t)
		if err != nil {
			fatalExec("go install "+t.Repository, err)
		}
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Filename to read/write the spec data.
const specfile = "tools.json"

type spec struct {
	Tools []*tool
}

func (s spec) write() error {
	specfilePath := filepath.Join(baseDirPath, specfile)

	f, err := os.Create(specfilePath)
	if err != nil {
		return fmt.Errorf("unable to open %s: %s", specfile, err)
	}
	defer f.Close()

	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json spec: %s", err)
	}

	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("unable to write %s: %s", specfile, err)
	}

	return nil
}

func (s spec) find(t *tool) int {
	for i, tt := range s.Tools {
		if t.Repository == tt.Repository {
			return i
		}
	}
	return -1
}

func (s spec) cleanup() {
	var pkgs []string
	for _, t := range s.Tools {
		pkgs = append(pkgs, t.Repository)
	}
	clean(pkgs)
}

func read() (spec, error) {
	specfilePath := filepath.Join(baseDirPath, specfile)

	file, err := os.Open(specfilePath)
	if err != nil {
		return spec{}, fmt.Errorf("unable to open spec file at %s: %s", specfile, err)
	}
	defer file.Close()

	s := new(spec)
	err = json.NewDecoder(file).Decode(s)
	if err != nil {
		return spec{}, err
	}
	return *s, nil
}

func specExists() bool {
	specfilePath := filepath.Join(baseDirPath, specfile)

	_, err := os.Stat(specfilePath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fatal("unable to stat tools.json: %s", err)
	}
	return true
}

func writeBlankSpec() error {
	return spec{[]*tool{}}.write()
}

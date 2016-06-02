package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Filename to read/write the spec data.
const specfile = "tools.json"

type spec struct {
	Tools []tool
}

func (s spec) write() error {
	f, err := os.Create(specfile)
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

func (s spec) add(t tool) {
	if s.find(t) != -1 {
		log(t.Repository + " already installed (did you mean retool upgrade?)")
		return
	}

	s.Tools = append(s.Tools, t)
	err := s.write()
	if err != nil {
		fatal("unable to add "+t.Repository, err)
	}

	s.sync()
}

func (s spec) remove(t tool) {
	idx := s.find(t)
	if idx == -1 {
		fatal(t.Repository+" is not in tools.json", nil)
	}
	s.Tools = append(s.Tools[:idx], s.Tools[idx+1:]...)

	err := s.write()
	if err != nil {
		fatal("unable to remove "+t.Repository, err)
	}

	s.sync()
}

func (s spec) upgrade(t tool) {
	idx := s.find(t)
	if idx == -1 {
		log(t.Repository + " is not yet installed (did you mean retool add?)")
		return
	}

	s.Tools[idx].Commit = t.Commit

	err := s.write()
	if err != nil {
		fatal("unable to remove "+t.Repository, err)
	}

	s.sync()
}

func (s spec) find(t tool) int {
	for i, tt := range s.Tools {
		if t.Repository == tt.Repository {
			return i
		}
	}
	return -1
}

func (s spec) sync() {
	for _, t := range s.Tools {
		err := install(t)
		if err != nil {
			fatalExec("failed to sync "+t.Repository, err)
		}
	}
}

func read() (spec, error) {
	file, err := os.Open(specfile)
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
	_, err := os.Stat(specfile)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fatal("unable to stat tools.json: %s", err)
	}
	return true
}

func writeBlankSpec() error {
	return spec{[]tool{}}.write()
}

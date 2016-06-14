package main

import (
	"encoding/json"
	"os"
	"path"
)

const manifestFile = "manifest.json"

type manifest map[string]string

func getManifest() manifest {
	m := manifest{}

	file, err := os.Open(path.Join(tooldir, manifestFile))
	if err != nil {
		return m
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&m)
	return m
}

func (m manifest) write() {
	f, err := os.Create(path.Join(tooldir, manifestFile))
	if err != nil {
		return
	}
	defer f.Close()

	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return
	}

	f.Write(bytes)
}

func (m manifest) outOfDate(ts []*tool) bool {
	// Make a copy to check for elements in ts but not m
	m2 := make(map[string]string)
	for k, v := range m {
		m2[k] = v
	}

	for _, t := range ts {
		if v, ok := m[t.Repository]; !ok || v != t.Commit {
			return true
		}
		delete(m2, t.Repository)
	}

	if len(m2) != 0 {
		return true
	}

	return false
}

func (m manifest) replace(ts []tool) {
	for k := range m {
		delete(m, k)
	}
	for _, t := range ts {
		m[t.Repository] = t.Commit
	}
}

package main

import (
	"os"
	"path/filepath"
)

func (s spec) sync() {
	m := getManifest()
	if m.outOfDate(s.Tools) {
		// Delete existing tools directory
		err := os.RemoveAll(toolDirPath)
		if err != nil {
			fatalExec("failed to remove _tools ", err)
		}

		// Recreate the tools directory
		err = ensureTooldir()
		if err != nil {
			fatal("failed to ensure tool dir", err)
		}

		// Download everything to tool directory
		for _, t := range s.Tools {
			err = download(t)
			if err != nil {
				fatalExec("failed to sync "+t.Repository, err)
			}
		}

		// Install the packages
		s.build()

		// Delete unneccessary source files
		s.cleanup()
		return
	}

	// code of tools has been cached, now check if binaries also cached
	binPath := filepath.Join(toolDirPath, "bin")
	binPathManifest := filepath.Join(binPath, manifestFile)
	buildBin := func() {
		s.build()
		s.cleanup()
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		buildBin()
		return
	}
	if _, err := os.Stat(binPathManifest); os.IsNotExist(err) {
		buildBin()
		return
	}

	m2 := getBinPathManifest()
	if !m.equals(m2) {
		buildBin()
	}
}

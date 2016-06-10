package main

import "os/exec"

func (s spec) sync() {
	m := getManifest()
	if m.outOfDate(s.Tools) {
		log("syncing")

		// Delete existing tools directory
		cmd := exec.Command("rm", "-r", "-f", tooldir)
		_, err := cmd.Output()
		if err != nil {
			fatalExec("failed to remove _tools ", err)
		}

		// Recreate the tools directory
		ensureTooldir()

		// Re-install everything
		for _, t := range s.Tools {
			err := install(t)
			if err != nil {
				fatalExec("failed to sync "+t.Repository, err)
			}
		}

		// Write a fresh manifest
		m.replace(s.Tools)
		m.write()

		// Delete unneccessary source files
		s.cleanup()
	}
}

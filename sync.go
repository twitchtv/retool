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

		// Download everything to cache
		for _, t := range s.Tools {
			err := download(t)
			if err != nil {
				fatalExec("failed to sync "+t.Repository, err)
			}
		}

		// Copy the cache into the tools directory
		cmd = exec.Command("cp", "-R", cacheDir+"/src", tooldir+"/src")
		_, err = cmd.Output()
		if err != nil {
			fatalExec("failed to copy data from cache ", err)
		}

		// Install the packages
		for _, t := range s.Tools {
			err = install(t)
			if err != nil {
				fatalExec("go install "+t.Repository, err)
			}
		}

		// Write a fresh manifest
		m.replace(s.Tools)
		m.write()

		// Delete unneccessary source files
		s.cleanup()
	}
}

package main

import (
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

type syncError struct {
	cmdName string
	err     error
}

func (e *syncError) Error() string {
	return fmt.Sprintf("failed to sync %s: %s", e.cmdName, e.err)
}

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
		var eg errgroup.Group
		for _, t := range s.Tools {
			t := t
			eg.Go(func() error {
				if err := download(t); err != nil {
					return &syncError{
						cmdName: t.Repository,
						err:     err,
					}
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			serr := err.(*syncError)
			fatalExec(serr.cmdName, serr.err)
		}

		// Install the packages
		s.build()

		// Delete unneccessary source files
		s.cleanup()
	}
}

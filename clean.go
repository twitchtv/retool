package main

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func clean(pkgs []string) {
	// This is rude but shouldn't be an issue...
	build.Default.GOPATH = tooldir

	// A recursive helper to take a list of packages and find their dependencies deeply
	found := map[string]struct{}{}
	var resolve func(string, []string) []string
	resolve = func(parent string, pkgs []string) []string {
		var r []string
		for _, pkg := range pkgs {
			if !strings.Contains(pkg, ".") {
				continue
			}

			p, err := build.Default.Import(pkg, path.Join(tooldir, "src", parent), 0)
			if err != nil {
				fatal(fmt.Sprintf("couldn't import package %q", pkg), err)
			}

			if _, ok := found[p.ImportPath]; ok {
				continue
			}

			found[p.ImportPath] = struct{}{}
			r = append(r, p.ImportPath)
			r = append(r, resolve(p.ImportPath, p.Imports)...)
		}
		return r
	}

	keep := resolve("", pkgs)
	base := path.Join(tooldir, "src")

	var toDelete []string
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// Bubble up errors
		if err != nil {
			return err
		}

		// Skip the root directory
		if base == path {
			return nil
		}

		// Get the package directory
		pkg, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		// Delete files in packages that aren't marked as "keep",
		// and any non-go or test files.
		if !info.IsDir() {
			pkg = filepath.Dir(pkg)
			_, keptPkg := found[pkg]
			isGo := strings.HasSuffix(path, ".go")
			isTest := strings.HasSuffix(path, "_test.go")
			if !keptPkg || !isGo || isTest {
				toDelete = append(toDelete, path)
			}
			return nil
		}

		// If the folder is a kept package or a parent, don't delete it and keep recursing
		for _, p := range keep {
			if strings.HasPrefix(p, pkg) {
				return nil
			}
		}

		// Otherwise this is a package that isn't imported at all. Delete it and stop recursing
		toDelete = append(toDelete, path)
		return filepath.SkipDir
	})
	if err != nil {
		fatal("unable to clean _tools", err)
	}

	for _, path := range toDelete {
		err = os.RemoveAll(path)
		if err != nil {
			fatal("unable to remove file or directory", err)
		}
	}
}

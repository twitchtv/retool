package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const version = "v1.0.1"

var cacheDir = ""

func init() {
	u, err := user.Current()
	if err == nil && u.HomeDir != "" {
		cacheDir = filepath.Join(u.HomeDir, ".retool")
	} else {
		cwd, err := os.Getwd()
		if err == nil {
			cacheDir = filepath.Join(cwd, ".retool")
		} else {
			cacheDir = ".retool"
		}
	}
}

func main() {
	flag.Parse()
	if err := ensureTooldir(); err != nil {
		fatal("failed to locate or create tool directory", err)
	}
	cmd, tool := parseArgs()

	if cmd == "version" {
		fmt.Fprintf(os.Stdout, "retool %s", version)
		os.Exit(0)
	}

	if !specExists() {
		if cmd == "add" {
			err := writeBlankSpec()
			if err != nil {
				fatal("failed to write blank spec", err)
			}
		} else {
			fatal("tools.json does not yet exist. You need to add a tool first with 'retool add'", nil)
		}
	}

	s, err := read()
	if err != nil {
		fatal("failed to load tools.json", err)
	}

	switch cmd {
	case "add":
		s.add(tool)
	case "upgrade":
		s.upgrade(tool)
	case "remove":
		s.remove(tool)
	case "build":
		s.build()
	case "sync":
		s.sync()
	case "do":
		s.sync()
		do()
	case "clean":
		err = os.RemoveAll(cacheDir)
		if err != nil {
			fatal("Failure during clean", err)
		}
	default:
		fatal(fmt.Sprintf("unknown cmd %q", cmd), nil)
	}
}

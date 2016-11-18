package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

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

	if !specExists() {
		if cmd == "add" {
			writeBlankSpec()
		} else {
			fatal("tools.json does not yet exist. You need to add a tool first with 'retool add'", nil)
		}
	}

	spec, err := read()
	if err != nil {
		fatal("failed to load tools.json", err)
	}

	switch cmd {
	case "add":
		spec.add(tool)
	case "upgrade":
		spec.upgrade(tool)
	case "remove":
		spec.remove(tool)
	case "sync":
		spec.sync()
	case "do":
		spec.sync()
		do()
	case "clean":
		os.RemoveAll(cacheDir)
	default:
		fatal(fmt.Sprintf("unknown cmd %q", cmd), nil)
	}
}

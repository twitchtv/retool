package main

import (
	"flag"
	"fmt"
	"os"
)

var verbose = flag.Bool("v", false, "verbose mode")

func init() {
	flag.Usage = func() {
		printUsageAndExit("", 0)
	}
	flag.Parse()
}

func printUsageAndExit(command string, exitCode int) {
	switch command {
	case "add":
		fmt.Println(addUsage)
	case "remove":
		fmt.Println(removeUsage)
	case "upgrade":
		fmt.Println(upgradeUsage)
	case "sync":
		fmt.Println(syncUsage)
	case "generate":
		fmt.Println(generateUsage)
	default:
		fmt.Println(usage)
	}
	os.Exit(exitCode)
}

func assertArgLength(command string, arglength int) {
	if len(os.Args) != arglength {
		printUsageAndExit(command, 1)
	}
}

var usage = `usage: retool [-v] (add | remove | upgrade | sync | generate | help)

use retool with a subcommand:

add will add a tool
remove will remove a tool
upgrade will upgrade a tool
sync will synchronize your tools with tools.json
generate will call 'go generate ./...' using your installed tools

help [command] will describe a command in more detail

For all commands, passing -v will enable verbose mode.
`

var addUsage = `usage: retool add [repository] [commit]

eg: retool add github.com/tools/godep 3020345802e4bff23902cfc1d19e90a79fae714e

Add will mark a repository as a tool you want to use. It will rewrite
tools.json to record this fact. It will then fetch the repository,
reset it to the desired commit, and install it to _tools/bin.

Don't use 'master' for the commit. It kind of defeats the whole purpose.
`

var upgradeUsage = `usage: retool upgrade [repository] [commit]

eg: retool upgrade github.com/tools/godep 3020345802e4bff23902cfc1d19e90a79fae714e

Upgrade set the commit SHA of a tool you want to use. It will
rewrite tools.json to record this fact. It will then fetch the
repository, reset it to the desired commit, and install it to
_tools/bin.

Don't use 'master' for the commit. It kind of defeats the whole purpose.
`

var removeUsage = `usage: retool remove [repository]

eg: retool remove github.com/tools/godep

Remove will remove a tool from your tools.json. It won't delete the
underlying repo from _tools, because it might be a dependency of some
other tool. If you really want to clean things up, you can nuke _tools
and it by calling 'rm -rf _tools && retool sync'.
`

var syncUsage = `usage: retool sync

Sync will synchronize your _tools directory to match tools.json.
`

var generateUsage = `usage: retool generate

retool generate will make sure your _tools directory is synced, and
then execute go generate ./... with the tools installed in _tools.

This is just
  retool sync && PATH=$PWD/_tools/bin:$PATH go generate ./...
That works too.
`

func parseArgs() (command string, t tool) {
	if len(os.Args) < 2 {
		printUsageAndExit("", 1)
	}

	command = os.Args[1]
	switch command {
	case "sync":
		assertArgLength(command, 2)
		return "sync", t

	case "add":
		assertArgLength(command, 4)
		t.Repository = os.Args[2]
		t.Commit = os.Args[3]
		return "add", t

	case "upgrade":
		assertArgLength(command, 4)
		t.Repository = os.Args[2]
		t.Commit = os.Args[3]
		return "upgrade", t

	case "remove":
		assertArgLength(command, 3)
		t.Repository = os.Args[2]

	case "generate":
		assertArgLength(command, 2)
		return "do", t

	case "help":
		assertArgLength(command, 3)
		printUsageAndExit(os.Args[2], 0)

	default:
		printUsageAndExit("", 1)
	}
	return "", t
}

func main() {
	ensureTooldir()
	cmd, tool := parseArgs()

	if !specExists(specfile) {
		if cmd == "add" {
			initializeSpec(specfile)
		} else {
			fatal("tools.json does not yet exist. You need to add a tool first with 'retool add'", nil)
		}
	}

	spec, err := read(specfile)
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
	case "generate":
		spec.sync()
		generate()
	}
}

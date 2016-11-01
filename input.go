package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	verboseFlag = flag.Bool("verbose", false, "Enable more detailed output that may be helpful for troubleshooting.")
	forkFlag    = flag.String("f", "", "Use a fork of the repository rather than the default upstream")

	// TODO: Refactor so that this global state is not necessary.
	positionalArgs []string
)

// TODO: Will this bother errcheck?
func verbosef(format string, a ...interface{}) (n int, err error) {
	if *verboseFlag {
		return fmt.Fprintf(os.Stderr, format, a...)
	}
	return 0, nil
}

func parseArgs() (command string, t *tool) {
	if !flag.Parsed() {
		panic("parseArgs expects that flags have already been parsed")
	}

	args := flag.Args()

	if len(args) < 1 {
		printUsageAndExit("", 1)
	}

	command = args[0]
	args = args[1:]
	t = new(tool)
	positionalArgs = args

	switch command {
	case "sync":
		assertArgLength(args, command, 0)
		return "sync", t

	case "add":
		assertArgLength(args, command, 2)
		t.Repository = args[0]
		t.ref = args[1]
		t.Fork = *forkFlag
		return "add", t

	case "upgrade":
		assertArgLength(args, command, 2)
		t.Repository = args[0]
		t.ref = args[1]
		t.Fork = *forkFlag
		return "upgrade", t

	case "remove":
		assertArgLength(args, command, 1)
		t.Repository = args[0]
		return "remove", t

	case "do":
		// A variable number of arguments are permissible for the 'do' subcommand; they are passed via t.PositionalArgs.
		return "do", t

	case "clean":
		assertArgLength(args, command, 0)
		return "clean", t

	case "help":
		assertArgLength(args, command, 1)
		printUsageAndExit(args[1], 0)

	default:
		printUsageAndExit("", 1)
	}
	return "", t
}

func assertArgLength(args []string, command string, arglength int) {
	if len(args) != arglength {
		printUsageAndExit(command, 1)
	}
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
	case "do":
		fmt.Println(doUsage)
	case "clean":
		fmt.Println(cleanUsage)
	default:
		fmt.Println(usage)
	}
	os.Exit(exitCode)
}

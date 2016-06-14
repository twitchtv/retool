package main

import (
	"fmt"
	"os"
)

func parseArgs() (command string, t *tool) {
	if len(os.Args) < 2 {
		printUsageAndExit("", 1)
	}

	command = os.Args[1]
	t = new(tool)

	switch command {
	case "sync":
		assertArgLength(command, 2)
		return "sync", t

	case "add":
		assertArgLength(command, 4)
		t.Repository = os.Args[2]
		t.ref = os.Args[3]
		return "add", t

	case "upgrade":
		assertArgLength(command, 4)
		t.Repository = os.Args[2]
		t.ref = os.Args[3]
		return "upgrade", t

	case "remove":
		assertArgLength(command, 3)
		t.Repository = os.Args[2]
		return "remove", t

	case "do":
		return "do", t

	case "help":
		assertArgLength(command, 3)
		printUsageAndExit(os.Args[2], 0)

	default:
		printUsageAndExit("", 1)
	}
	return "", t
}

func assertArgLength(command string, arglength int) {
	if len(os.Args) != arglength {
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
	default:
		fmt.Println(usage)
	}
	os.Exit(exitCode)
}

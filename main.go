package main

func main() {
	ensureTooldir()
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
	}
}

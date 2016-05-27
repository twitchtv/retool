package main

var usage = `usage: retool (add | remove | upgrade | sync | generate | help)

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

var doUsage = `usage: retool do [command and args]

retool do will make sure your _tools directory is synced, and then
execute a command with the tools installed in _tools.

This is just
  retool sync && PATH=$PWD/_tools/bin:$PATH [command and args]
That works too.
`

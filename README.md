# Retool: Make Code Generation Great Again #

## hao git ##

```sh
go get github.com/golang/protobuf/protoc-gen-go
```

## wat do ##

Add a tool dependency:
```sh
retool add github.com/golang/protobuf/protoc-gen-go origin/master
```

Use it to generate code:
```sh
# calls go generate ./... using only tools installed with retool
retool do go generate ./...
```

Upgrade your tools to their latest version:
```sh
retool upgrade github.com/golang/protobuf/protoc-gen-go origin/master
```

Stop using that stupid tool you dont like anymore:
```sh
retool remove github.com/tools/godep
```

Stay in sync:
```sh
# makes sure your tools match tools.json
retool sync
```

## what is this ##

retool helps manage the versions of _tools_ that you use with your
repository. These are executables that are a crucial part of your
development environment, but aren't imported by any of your code, so
they don't get scooped up by godeps (or any other vendoring tool).

Some examples of tools:

 - [github.com/tools/godep](github.com/tools/godep) is a tool to
   vendor Go packages.
 - [github.com/golang/protobuf/protoc-gen-go](https://github.com/golang/protobuf/protoc-gen-go)
   is a tool to compile Go code from protobuf definitions.
 - [github.com/maxbrunsfeld/counterfeiter](https://github.com/maxbrunsfeld/counterfeiter)
   is a tool to generate mocks of interfaces.

## why would i need to manage these things  ##

**TL;DR:** if you work with anyone else on your project, and they have
different versions of their tools, everything turns to shit.

One of the best parts about Go is that it is very, very simple. This
makes it straightforward to write code generation utilities. You don't
need to generate code for every project, but in large ones, code
generation can help you be much more productive.

Like, if you're writing tests that use an interface, you can use code
generation to quickly whip up structs which mock the interface so you
can force them to return errors. This way, you can test edge cases for
your interaction points with
interfaces. [github.com/maxbrunsfeld/counterfeiter](https://github.com/maxbrunsfeld/counterfeiter)
does this pretty well!

If you want to use the generate code, you should check in the
generated `.go` code to git, not just the sources, so that build boxes
and the like don't need all these code generation tools, and so that
`go get` just works cleanly.

This poses a problem, though, as soon as you start working with other
people on your project: if you have different versions of your code
generation tools, which generate slightly different output, you'll get
lots of meaningless churn in your commits. This sucks! There has to be
a better way!

## the retool way ##

retool records the versions of tools you want in a file,
`tools.json`. The file looks like this:

```json
{
  "Tools": [
    {
      "Repository": "code.justin.tv/common/twirp/protoc-gen-twirp",
      "Commit": "24eb49ba0f7cd692f60f11af1cba4a515ab64e06"
    },
    {
      "Repository": "github.com/golang/protobuf/protoc-gen-go",
      "Commit": "2fea9e168bab814ca0c6e292a6be164f624fc6ca"
    }
  ]
}
```

Tools are identified by repo and commit. Each tool in `tools.json`
will be installed to `_tools`, which is a private GOPATH just
dedicated to keeping track of these tools.

In practice, you don't need to know much about `tools.json`. You check
it in to git so that everybody stays in sync, but you manage it with
`retool add|upgrade|remove`.

When it's time to generate code, **instead of `go generate ./...`**,
you use `retool do go generate ./...` to use your sweet, vendored
tools. This really just calls `PATH=$PWD/_tools/bin:PATH go generate
./...`; if you want to do anything fancy, you can feel free to use
that path too.

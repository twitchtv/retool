# go-gh-releaser #

This is a little utility for cross-compiling and releasing Go projects with
github releases. It cross-compiles a project using
[gox](https://github.com/mitchellh/gox) and then uploads the resulting bins to
the right release. Auth is handled with a personal access token.

## Installation ##

```
go get github.com/mitchellh/gox
go get github.com/spenczar/go-gh-releaser
```

## Usage ##
Use it thusly, subsitituting your real github access token:

```
go-gh-releaser -release <RELEASE-TAG> -repo <OWNER/REPO> -token <YOUR-ACCESS-TOKEN>
```

For example:
```
go-gh-releaser -release v1.3.5 -repo twitchtv/retool -token dc------------------------------------84
```

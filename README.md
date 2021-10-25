# forder

Forder check order of Go functions.

## Installation

To install the `forder` command, run

```
$ go install github.com/hawkingrei/forder@latest
```

and put the resulting binary in one of your PATH directories if
`$GOPATH/bin` isn't already in your PATH.

## Usage

```
Check order of Go functions.
Usage:
    forder [flags] <Go file or directory> ...
Flags:
	-ignore REGEX         exclude files matching the given regular expression
	-not-skip-vender      not skip vendor folders
	-skip-test            skip test file
```
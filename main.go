package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

const usageDoc = `Check order of Go functions.
Usage:
    forder [flags] <Go file or directory> ...
Flags:
	-ignore REGEX         exclude files matching the given regular expression
	-not-skip-vender      not skip vendor folders
	-skip-test            skip test file
`

func main() {
	ignore := flag.String("ignore", "", "exclude files matching the given regular expression")
	notSkipVendor := flag.Bool("not-skip-vender", false, "not skip vendor folders")
	skipTest := flag.Bool("skip-test", false, "skip test file")
	log.SetFlags(0)
	log.SetPrefix("forder: ")
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		usage()
	}

	rg := regex(*ignore)
	config := Config{
		SkipExp:    rg,
		SkipVendor: !*notSkipVendor,
		SkipTest:   *skipTest,
		Paths:      paths,
	}
	Analyze(config)
}

func usage() {
	_, _ = fmt.Fprint(os.Stderr, usageDoc)
	os.Exit(2)
}

func regex(expr string) *regexp.Regexp {
	if expr == "" {
		return nil
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		log.Fatal(err)
	}
	return re
}

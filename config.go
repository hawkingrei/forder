package main

import "regexp"

type Config struct {
	SkipVendor bool
	SkipTest   bool
	SkipExp    *regexp.Regexp
	Paths      []string
}

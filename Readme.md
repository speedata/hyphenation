[![GoDoc](https://godoc.org/github.com/speedata/hyphenation?status.svg)](https://godoc.org/github.com/speedata/hyphenation) [![CircleCI](https://circleci.com/gh/speedata/hyphenation/tree/master.svg?style=shield)](https://circleci.com/gh/speedata/hyphenation/tree/master)

A port of TeX's hyphenation algorithm to Go
===========================================

Installation
------------

    go get github.com/speedata/hyphenation


Prerequisites
-------------

Download a hyphenation pattern file from CTAN, for example from <https://ctan.math.utah.edu/ctan/tex-archive/language/hyph-utf8/tex/generic/hyph-utf8/patterns/txt/>

Usage
-----

````go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/speedata/hyphenation"
)

func main() {
	filename := "hyph-en-us.pat.txt"
	r, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	l, err := hyphenation.New(r)
	if err != nil {
		log.Fatal(err)
	}

	var h []int
	for _, v := range []string{"Computer", "developers"} {
		h = l.Hyphenate(v)
		fmt.Println(v, h)
	}
}
````

Other
------
Contact: <gundlach@speedata.de><br>
Twitter: [@speedata](https://twitter.com/speedata)<br>
License: cc0 / public domain (<https://creativecommons.org/publicdomain/zero/1.0/>)<br>
Status: Just an example hack, never used it in production yet

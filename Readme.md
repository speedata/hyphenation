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
		fmt.Println(v, h) // [3 6] and [2 5 7 9]
	}
}
````

Debugging hyphenation patterns
==============================

Similar to getting the hyphenation slice, you can get a detailed view of the hyphenation patterns used in a word:

````go
str := l.DebugHyphenate("developers")
fmt.Println(str)
````


results in

       .   d   e   v   e   l   o   p   e   r   s   .
         0   0   1   0   |   |   |   |   |   |   |    de1v
         |   0   0   0   0   3   0   |   |   |   |    evel3o
         |   |   0   0   4   0   0   |   |   |   |    ve4lo
         |   |   |   |   |   0   0   1   0   0   |    op1er
         |   |   |   |   |   |   |   0   0   1   0    er1s
         |   |   |   |   |   |   |   |   4   0   2    4rs2
    max: 0   0   1   0   4   3   0   1   4   1   2
    final: d   e - v   e   l - o   p - e   r - s

Other
------
Contact: <gundlach@speedata.de><br>
Twitter: [@speedata](https://twitter.com/speedata)<br>
License: cc0 / public domain (<https://creativecommons.org/publicdomain/zero/1.0/>)<br>
Status: Just an example hack, never used it in production yet

[![Go Reference](https://pkg.go.dev/badge/github.com/speedata/hyphenation.svg)](https://pkg.go.dev/github.com/speedata/hyphenation)&nbsp;[![Explore in Constellation](https://img.shields.io/badge/Explore%20in-Constellation-blue)](https://constellation.speedata.de/#hyphenation)

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

Sharing patterns
----------------

A `Lang` keeps every word it hyphenates in a cache. To parse a pattern file once and use it for many documents, give each document a clone: it shares the patterns and has its own cache and minimums, which go away with it.

````go
doc := l.Clone()
doc.Leftmin, doc.Rightmin = 2, 3
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
Mastodon: [@speedata@typo.social](https://typo.social/@speedata)<br>
License: cc0 / public domain (<https://creativecommons.org/publicdomain/zero/1.0/>)<br>
Status: used in the [boxes and glue](https://github.com/boxesandglue/boxesandglue) library.

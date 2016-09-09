// Public domain.
// Comments to gundlach@speedata.de

// Package hyphenation hyphenates words with TeXs algorithm.
//
// The algorithm is used in TeX and originally from Franklin Liang, his thesis can be downloaded from
// https://www.tug.org/docs/liang%20/liang-thesis.pdf
//
// You need pattern files which can be downloaded from
// http://ctan.math.utah.edu/ctan/tex-archive/language/hyph-utf8/tex/generic/hyph-utf8/patterns/txt/
package hyphenation

import (
	"bufio"
	"io"
	"unicode"
)

// Lang is a language object for hyphenation. Use it by calling New(), otherwise the object is not initialized properly.
type Lang struct {
	patterns map[string][]byte
}

// New loads patterns from the reader. Patterns are word substrings with a hyphenation priority
// between each letter, 0s omitted. Example patterns are “.ach4 at3est 4if.” where a dot
// denotes a word boundary. An odd number means “don't hyphenate here”, everything else allows
// hyphenation at this point. The final priority for each position is the maximum of each priority
// given in each applied pattern.
func New(r io.Reader) (*Lang, error) {
	l := &Lang{}
	l.patterns = make(map[string][]byte)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	var wordpart []rune
	var pat []byte

	var prio byte

	// for each word in the pattern file
	for s.Scan() {
		pattern := s.Text()
		wordpart = wordpart[:0]
		pat = []byte{}
		var prev rune

		for _, v := range pattern {
			// if it is a letter, the priority is the previous
			// entry (if it is a digit) or 0
			if !unicode.IsDigit(v) {
				if unicode.IsDigit(prev) {
					prio = byte(prev) - '0'
				} else {
					prio = 0
				}
				wordpart = append(wordpart, v)
				if v != '.' {
					pat = append(pat, prio)
				}
			}
			prev = v
		}
		if unicode.IsDigit(prev) {
			pat = append(pat, byte(prev)-'0')
		} else {
			pat = append(pat, 0)
		}

		l.patterns[string(wordpart)] = pat
		// pattern | wordpart | pat
		// uto1    | uto      | 0,0,0,1
		// .un3g   | .ung     | 0,0,3,0
		// h2en.   | hen.     | 0,2,0,0
	}
	return l, s.Err()
}

// Hyphenate returns an array of int with resulting break points.
// For example the word “developers” with English (US) hyphenation
// patterns could return [2 5 7 9] which means de-vel-op-ers
func (l *Lang) Hyphenate(word string) []int {
	var rword []rune
	for _, letter := range "." + word + "." {
		rword = append(rword, unicode.ToLower(letter))
	}
	var startpos int
	var wordpart []rune
	// maxPrio[i] contains the maximum priority before the letter i
	maxPrio := make([]byte, len(rword)-1)

	// generate all possible substrings for the word
	for j := 1; j < len(rword); j++ {
		for i := j + 1; i <= len(rword); i++ {

			startpos = j - 1
			wordpart = rword[j-1 : i]

			// if there is a pattern for this substring
			if pattern, ok := l.patterns[string(wordpart)]; ok {
				// pattern is a byte slice such as [0 0 1 0]
				// associated with the word part, "dev" for example
				// the pattern in the hyphenation file would be
				// de1v or similar 0d0e1v0
				// when the pattern contains a dot at the beginning, this
				// marks the start of the word. Therefore there is no
				// priority before the .
				if wordpart[0] == '.' {
					startpos = startpos + 1
				}
				for i := 0; i < len(pattern); i++ {
					if pattern[i] > maxPrio[startpos-1+1] {
						maxPrio[startpos-1+i] = pattern[i]
					}
				}
			}
		}
	}

	// the odd entries in maxPrio are valid break points
	var positions []int
	for i := 1; i < len(maxPrio); i++ {
		if maxPrio[i]%2 != 0 {
			positions = append(positions, i)
		}
	}
	return positions
}

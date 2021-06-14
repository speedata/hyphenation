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
	"strings"
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

type patternposition struct {
	startpos int
	wordpart []rune
	pattern  []byte
}

func (l *Lang) doHyphenate(rword []rune) []patternposition {
	var patterninfo []patternposition
	var startpos int
	var wordpart []rune
	// generate all possible substrings for the word
	for j := 1; j < len(rword); j++ {
		startpos = j - 1
		// if rword[0] == '.' {
		// 	startpos = startpos + 1
		// }
		for i := j + 1; i <= len(rword); i++ {
			wordpart = rword[startpos:i]

			// if there is a pattern for this substring
			if pattern, ok := l.patterns[string(wordpart)]; ok {
				patterninfo = append(patterninfo, patternposition{
					startpos: startpos,
					wordpart: wordpart,
					pattern:  pattern,
				})
			}
		}
	}
	return patterninfo
}

// DebugHyphenate returns a multi-line string with information about the patterns used and the priorities.
func (l *Lang) DebugHyphenate(word string) string {
	var rword []rune
	for _, letter := range "." + word + "." {
		rword = append(rword, unicode.ToLower(letter))
	}
	breakpoints := l.doHyphenate(rword)
	maxPrio := make([]byte, len(rword)-1)
	var b strings.Builder
	b.WriteString("   ")
	for _, v := range rword {
		b.WriteRune(v)
		b.WriteString("   ")
	}
	b.WriteByte(10)

	for _, bp := range breakpoints {
		startpos := bp.startpos
		sub := 1

		if bp.wordpart[0] == '.' {
			sub = 0
		}
		for i := 0; i < len(bp.pattern); i++ {
			if bp.pattern[i] > maxPrio[startpos-sub+i] {
				maxPrio[startpos-sub+i] = bp.pattern[i]
			}
		}

		b.WriteString("    ")
		for i := 0; i < bp.startpos-sub; i++ {
			b.WriteString(" |  ")
		}
		for i := 0; i < len(bp.pattern); i++ {
			b.WriteRune(' ')
			b.WriteRune('0' + rune(bp.pattern[i]))
			b.WriteRune(' ')
			b.WriteRune(' ')
		}
		for i := bp.startpos + len(bp.pattern) - sub; i < len(rword); i++ {
			b.WriteString(" |  ")
		}
		b.WriteRune(' ')
		for i := 0; i < len(bp.wordpart); i++ {
			if bp.pattern[i] > 0 {
				b.WriteRune(rune(bp.pattern[i]) + '0')
			}
			b.WriteRune(bp.wordpart[i])
		}
		lastprio := bp.pattern[len(bp.pattern)-1]
		if lastprio > 0 {
			b.WriteRune(rune(lastprio) + '0')
		}
		b.WriteByte(10)
	}
	b.WriteString("max: ")
	for i := 0; i < len(maxPrio); i++ {
		b.WriteRune('0' + rune(maxPrio[i]))
		b.WriteString("   ")
	}
	b.WriteByte(10)

	b.WriteString("final: ")
	for i := 1; i < len(rword)-1; i++ {
		b.WriteRune(rword[i])
		b.WriteRune(' ')
		if maxPrio[i]%2 == 0 {
			b.WriteRune(' ')
		} else {
			b.WriteRune('-')
		}
		b.WriteRune(' ')

	}

	b.WriteByte(10)
	return b.String()
}

// Hyphenate returns an array of int with resulting break points.
// For example the word “developers” with English (US) hyphenation
// patterns could return [2 5 7 9] which means de-vel-op-er-s
func (l *Lang) Hyphenate(word string) []int {
	var rword []rune
	for _, letter := range "." + word + "." {
		rword = append(rword, unicode.ToLower(letter))
	}

	breakpoints := l.doHyphenate(rword)
	maxPrio := make([]byte, len(rword)-1)

	for _, bp := range breakpoints {
		// pattern is a byte slice such as [0 0 1 0]
		// associated with the word part, "dev" for example
		// the pattern in the hyphenation file would be
		// de1v or similar 0d0e1v0
		// when the pattern contains a dot at the beginning, this
		// marks the start of the word. Therefore there is no
		// priority before the .
		startpos := bp.startpos
		if bp.wordpart[0] == '.' {
			startpos++
		}
		for i := 0; i < len(bp.pattern); i++ {
			if bp.pattern[i] > maxPrio[startpos-1+i] {
				maxPrio[startpos-1+i] = bp.pattern[i]
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

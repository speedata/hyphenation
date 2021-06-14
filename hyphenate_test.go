package hyphenation

import (
	"strings"
	"testing"
)

const patternsEN string = `1co
4m1p
pu2t
5pute
put3er
pos1s
1pos
2ess
2ss
s1e4s
s1si
1sio
5sion
2io
o2n
`

const patternsDE string = `
.sc4h4
1de
1fa
1fe
1ne
1sc
1se
1ße
1ta
2ä3s2e
2hn
2hr
2if
4ahr
4f1f
4l1l
4n1g
4n1n
4r1t
4rn
5fahrt
6r1s
c2h
g2r4
gl2
h2ü
h3ner
hn2e
hühne4
ng2läs
ö1ß
rn3g2
s1t
s2ta
s4er.
üh3ne
`

func TestHyphenate(t *testing.T) {
	r := strings.NewReader(patternsEN)
	l, err := New(r)
	if err != nil {
		t.Error(err)
	}

	data := []struct {
		word        string
		breakpoints []int
	}{
		{"possession", []int{3, 6}},
		{"Computer", []int{3, 6}},
		{"a", []int{}},
	}

	for _, entry := range data {
		h := l.Hyphenate(entry.word)
		if len(h) != len(entry.breakpoints) {
			t.Errorf("Hyphenate(%s), len = %d, want %d", entry.word, len(h), len(entry.breakpoints))
		}

		for i, v := range h {
			if bp := entry.breakpoints[i]; bp != v {
				t.Errorf("Hyphenate(%s), breakpoint[%d] = %d, want %d", entry.word, i, bp, v)
			}
		}
	}
}

func TestHyphenateDE(t *testing.T) {
	r := strings.NewReader(patternsDE)
	l, err := New(r)
	if err != nil {
		t.Error(err)
	}

	data := []struct {
		word        string
		breakpoints []int
	}{
		{"größer", []int{3}},
		{"Schiffahrt", []int{5, 9}},
		{"Hühnerstall", []int{3, 6, 10}},
		{"ferngläser", []int{4, 7}},
		{"denn", []int{3}},
	}

	for _, entry := range data {
		h := l.Hyphenate(entry.word)
		if len(h) != len(entry.breakpoints) {
			t.Errorf("Hyphenate(%s), len = %d, want %d", entry.word, len(h), len(entry.breakpoints))
		}

		for i, v := range h {
			if bp := entry.breakpoints[i]; bp != v {
				t.Errorf("Hyphenate(%s), breakpoint[%d] = %d, want %d", entry.word, i, bp, v)
			}
		}
	}
}

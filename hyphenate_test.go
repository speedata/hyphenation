package hyphenation

import (
	"slices"
	"strings"
	"sync"
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

func TestHyphenateDEmin(t *testing.T) {
	r := strings.NewReader(patternsDE)
	l, err := New(r)
	if err != nil {
		t.Error(err)
	}
	l.Leftmin = 2
	l.Rightmin = 3

	data := []struct {
		word        string
		breakpoints []int
	}{
		{"größer", []int{3}},
		{"Schiffahrt", []int{5}},
		{"Hühnerstall", []int{3, 6}},
		{"ferngläser", []int{4, 7}},
		{"denn", []int{}},
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

// One Lang hyphenates in several goroutines at once, each word new to the
// cache in one of them.
func TestHyphenateConcurrently(t *testing.T) {
	l, err := New(strings.NewReader(patternsEN))
	if err != nil {
		t.Fatal(err)
	}
	var words []string
	for _, a := range "abcdefghij" {
		for _, b := range "klmnopqrst" {
			words = append(words, "com"+string(a)+"puter"+string(b)+"possession")
		}
	}
	want := make([][]int, len(words))
	single, _ := New(strings.NewReader(patternsEN))
	for i, w := range words {
		want[i] = single.Hyphenate(w)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, w := range words {
				if got := l.Hyphenate(w); !slices.Equal(got, want[i]) {
					t.Errorf("Hyphenate(%s) = %v, want %v", w, got, want[i])
				}
			}
		}()
	}
	wg.Wait()
}

// A clone hyphenates as the original, with the original's minimums, shares
// its patterns and keeps a word cache of its own.
func TestClone(t *testing.T) {
	l, err := New(strings.NewReader(patternsDE))
	if err != nil {
		t.Fatal(err)
	}
	l.Leftmin, l.Rightmin = 2, 3
	c := l.Clone()
	if c.Leftmin != 2 || c.Rightmin != 3 {
		t.Errorf("clone minimums %d and %d, want 2 and 3", c.Leftmin, c.Rightmin)
	}
	for _, w := range []string{"größer", "Schiffahrt", "Hühnerstall", "ferngläser"} {
		if got, want := c.Hyphenate(w), l.Hyphenate(w); !slices.Equal(got, want) {
			t.Errorf("clone Hyphenate(%s) = %v, want %v", w, got, want)
		}
	}
	l.Hyphenate("Sonnenschein")
	if _, ok := c.hyphenatedWords[".sonnenschein."]; ok {
		t.Error("a word the original hyphenated is in the clone's cache")
	}
	if len(c.hyphenatedWords) != 4 || len(l.hyphenatedWords) != 5 {
		t.Errorf("caches hold %d (clone) and %d (original) words, want 4 and 5", len(c.hyphenatedWords), len(l.hyphenatedWords))
	}
	l.patterns["zzz"] = []byte{0, 0, 0, 0}
	if _, ok := c.patterns["zzz"]; !ok {
		t.Error("the clone does not share the original's patterns")
	}
}

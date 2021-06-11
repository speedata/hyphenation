package hyphenation

import (
	"strings"
	"testing"
)

const patterns string = `1co
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

func TestHyphenate(t *testing.T) {
	r := strings.NewReader(patterns)
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

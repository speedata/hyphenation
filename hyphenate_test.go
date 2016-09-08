package hyphenation

import (
	"strings"
	"testing"
)

func TestHyphenate(t *testing.T) {
	r := strings.NewReader("1co 4m1p pu2t 5pute put3er")
	l, err := New(r)
	if err != nil {
		t.Error(err)
	}

	h := l.Hyphenate("Computer")
	expected := []int{3, 6}
	if len(h) != len(expected) {
		t.Error("Hyphenation error, different length")
	}
	for i, v := range h {
		if expected[i] != v {
			t.Error("Hyphenation error")
		}
	}
}

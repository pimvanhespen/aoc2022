package main

import (
	"strings"
	"testing"
)

const testInput = `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`

func TestSolve1(t *testing.T) {
	v, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := solve1(v), 64; got != want {
		t.Errorf("solve1(%v) = %d, want %d", v, got, want)
	}
}

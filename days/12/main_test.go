package main

import (
	"testing"
)

var input = []byte(`Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi`)

func Test_solve1(t *testing.T) {

	const want = 31

	chart, b, e := parse(input)

	if got := solve1(chart, b, e); got != want {
		t.Errorf("solve1() = %v, want %v", got, want)
	}
}

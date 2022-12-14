package main

import (
	"strings"
	"testing"
)

const testInput = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

func TestSolve1(t *testing.T) {
	const want = 24
	field := parse(strings.NewReader(testInput))

	got := solve1(field)
	if got != want {
		t.Errorf("got != want: got %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	const want = 93
	field := parse(strings.NewReader(testInput))

	got := solve2(field)
	if got != want {
		t.Errorf("got != want: got %d, want %d", got, want)
	}
}

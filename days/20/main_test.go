package main

import (
	"strings"
	"testing"
)

var testInput = `1
2
-3
3
-2
0
4`

func TestSolve1(t *testing.T) {
	const want = 3
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	result := solve1(input)
	if result != want {
		t.Errorf("Expected %d, got %d", want, result)
	}
}

func TestSolve2(t *testing.T) {
	const want = 1623178306
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	result := solve2(input)
	if result != want {
		t.Errorf("Expected %d, got %d", want, result)
	}
}

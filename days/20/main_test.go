package main

import (
	"strings"
	"testing"
)

var testInput = ``

func TestSolve1(t *testing.T) {
	const want = 0
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	result := solve1(input)
	if result != 0 {
		t.Errorf("Expected %d, got %d", want, result)
	}
}

func TestSolve2(t *testing.T) {
	const want = 0
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	result := solve2(input)
	if result != 0 {
		t.Errorf("Expected %d, got %d", want, result)
	}
}

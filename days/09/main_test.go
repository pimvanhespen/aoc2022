package main

import (
	"strings"
	"testing"
)

func TestSolve1(t *testing.T) {
	reader := strings.NewReader(`R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`)
	in, err := parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	got := solve1(in)
	const want = 13
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	reader := strings.NewReader(`R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20`)
	in, err := parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	got := solve2(in)
	const want = 36
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

package main

import (
	"io"
	"strings"
	"testing"
)

var numbers = []string{
	"1",
	"2",
	"3",
	"4",
	" ",
	"1",
	"2",
	"3",
	" ",
	"1",
	"2",
	" ",
	"1",
}

type readCloser struct {
	io.Reader
}

func (r readCloser) Close() error {
	return nil
}

func Test_parseInput(t *testing.T) {
	text := strings.Join(numbers, "\n")
	r := strings.NewReader(text)
	rc := readCloser{r}

	invs, err := parseInput(rc)
	if err != nil {
		t.Fatal(err)
	}

	if len(invs) != 4 {
		t.Fatalf("expected 4 inventories, got %d", len(invs))
	}

	for i, inv := range invs {
		expected := len(invs) - i
		found := len(inv.items)

		if found != expected {
			t.Errorf("expected inventory %d to have %d items, got %d\n%v", i, expected, found, inv)
		}
	}
}

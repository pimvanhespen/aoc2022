package main

import (
	"strings"
	"testing"
)

var testIn = `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`

func TestSolve1(t *testing.T) {
	want := 95437
	root, err := parse(strings.NewReader(testIn))
	if err != nil {
		t.Fatal(err)
	}
	got := solve1(root)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

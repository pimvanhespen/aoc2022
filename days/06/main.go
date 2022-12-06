package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	rc, err := aoc.Get(6)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	bts, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part1", solve1(bts))
	fmt.Println("Part2", solve2(bts))
}

func solve1(bts []byte) int {
	for i := range bts {
		if allDifferent(bts[i : i+4]) {
			return i + 4
		}
	}
	return 0
}

func solve2(bts []byte) int {
	for i := range bts {
		if allDifferent(bts[i : i+14]) {
			return i + 14
		}
	}
	return 0
}

func allDifferent(bts []byte) bool {
	for i := 0; i < len(bts)-1; i++ {
		if bytes.IndexByte(bts[i+1:], bts[i]) != -1 {
			return false
		}
	}
	return true
}

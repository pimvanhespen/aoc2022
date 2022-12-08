package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	r, err := aoc.Get(6)
	if err != nil {
		panic(err)
	}

	bts, err := io.ReadAll(r)
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

// -- alternatives --

func solve_map(bts []byte, size int) int {
	var m = make(map[byte]int)
	var c byte
	for i := 0; i < len(bts); i++ {
		if i > size-1 {

			c = bts[i-size]

			n, ok := m[c]
			if !ok {
				panic("missing value in map")
			}

			if n == 1 {
				delete(m, c)
			} else {
				m[c]--
			}
		}

		c = bts[i]
		if n, ok := m[c]; ok {
			m[c] = n + 1
		} else {
			m[c] = 1
		}

		if len(m) == size {
			return i + 1
		}
	}
	return -1
}

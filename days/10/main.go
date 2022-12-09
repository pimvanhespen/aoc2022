package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	r, err := aoc.Get(10)
	if err != nil {
		panic(err)
	}

	in, err := parse(r)
	if err != nil {
		panic(err)
	}

	// Part 1
	fmt.Println("Part 1:", solve1(in))
	fmt.Println("Part 2:", solve2(in))
}

func parse(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}

func solve1(in []byte) int {
	return 0
}

func solve2(in []byte) int {
	return 0
}

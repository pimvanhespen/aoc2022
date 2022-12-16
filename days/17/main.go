package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	r, err := aoc.Get(17)
	if err != nil {
		panic(err)
	}

	input, err := parse(r)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))
}

type Input struct {
	// todo
}

func parse(reader io.Reader) (Input, error) {
	panic("not implemented")
}

func solve1(input Input) int {
	panic("not implemented")
}

func solve2(input Input) int {
	panic("not implemented")
}

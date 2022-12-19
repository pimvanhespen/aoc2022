package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	input, err := aoc.Load(20, parse)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))
}

type Input struct{}

func parse(reader io.Reader) (Input, error) {
	return Input{}, nil
}

func solve1(input Input) int {
	return 0
}

func solve2(input Input) int {
	return 0
}

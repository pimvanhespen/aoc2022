package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	reader, err := aoc.Get(9)
	if err != nil {
		panic(err)
	}

	input, err := parse(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))

}

type Input []byte

func parse(reader io.Reader) (Input, error) {

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return Input(b), nil
}

func solve1(input Input) int {
	return 0
}

func solve2(input Input) int {
	return 0
}
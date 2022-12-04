package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/list"
)

type Pair struct {
	Left, Right Range
}

type Range struct {
	Min, Max int
}

func (r Range) Contains(other Range) bool {
	return r.Min <= other.Min && r.Max >= other.Max
}

func (r Range) Overlaps(other Range) bool {
	return r.Min <= other.Max && r.Max >= other.Min
}

func main() {
	rc, err := aoc.Get(4)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	p := aoc.Parser[Pair]{
		SkipEmptyLines: false,
		ParseFn:        parseLine,
	}

	rows, err := p.Rows(rc)
	if err != nil {
		panic(err)
	}

	// Part 1
	engulfing := list.Count(rows, func(p Pair) bool {
		return p.Left.Contains(p.Right) || p.Right.Contains(p.Left)
	})
	fmt.Println("Part 1:", engulfing)

	// Part 2
	overlapping := list.Count(rows, func(p Pair) bool {
		return p.Left.Overlaps(p.Right)
	})
	fmt.Println("Part 2:", overlapping)
}

func parseLine(line string) (Pair, error) {
	const format = "%d-%d,%d-%d"

	var row Pair
	_, err := fmt.Sscanf(line, format, &row.Left.Min, &row.Left.Max, &row.Right.Min, &row.Right.Max)
	if err != nil {
		return Pair{}, err
	}
	return row, nil
}

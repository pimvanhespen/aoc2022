package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
)

type Row struct{}

func toRow(line string) (Row, error) {
	return Row{}, nil
}

func main() {

	rc, err := aoc.Get(5)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	parser := aoc.Parser[Row]{
		SkipEmptyLines: true,
		ParseFn:        toRow,
	}

	rows, err := parser.Rows(rc)
	if err != nil {
		panic(err)
	}

	fmt.Print("Part 1:")
	fmt.Println(solve1(rows))

	fmt.Print("Part 2:")
	fmt.Println(solve2(rows))
}

func solve1(rows []Row) (int, error) {
	return 0, nil
}

func solve2(rows []Row) (int, error) {
	return 0, nil
}

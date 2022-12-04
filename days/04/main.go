package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/puzzleinput"
	"io"
)

func main() {
	rc, err := puzzleinput.Get(4)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	input, err := parseInput(rc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))
}

func solve1(rows []Row) int {
	var count int
	for _, row := range rows {
		if row.Left.Contains(row.Right) || row.Right.Contains(row.Left) {
			count++
		}
	}
	return count
}

func solve2(rows []Row) int {
	var count int
	for _, row := range rows {
		if row.Left.Overlaps(row.Right) {
			count++
		}
	}
	return count
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

type Row struct {
	Left, Right Range
}

func parseInput(reader io.Reader) ([]Row, error) {
	scanner := bufio.NewScanner(reader)
	var rows []Row
	for scanner.Scan() {
		line := scanner.Text()
		row, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func parseLine(line string) (Row, error) {
	const format = "%d-%d,%d-%d"

	var row Row
	_, err := fmt.Sscanf(line, format, &row.Left.Min, &row.Left.Max, &row.Right.Min, &row.Right.Max)
	if err != nil {
		return Row{}, err
	}
	return row, nil
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/list"
	"github.com/pimvanhespen/aoc2022/pkg/set"
	"io"
)

func main() {
	in, err := aoc.Get(3)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	rows, err := parseInput(in)
	if err != nil {
		panic(err)
	}

	result, err := solve1(rows)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", result)

	result, err = solve2(rows)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result)
}

func solve1(rucksacks [][]byte) (int, error) {

	var total int

	for _, r := range rucksacks {
		half := len(r) / 2

		left := set.New(r[:half]...)
		right := set.New(r[half:]...)

		intersection := left.Intersection(right)

		if intersection.Len() != 1 {
			return 0, errors.New("intersection is not one element")
		}

		total += value(intersection.ToSlice()[0])
	}

	return total, nil
}

func solve2(rs [][]byte) (int, error) {
	if len(rs)%3 != 0 {
		panic("not a multiple of 3")
	}

	var total int

	for i := 0; i < len(rs); i += 3 {

		s1 := set.New(rs[i]...)
		s2 := set.New(rs[i+1]...)
		s3 := set.New(rs[i+2]...)

		badge := s1.Intersection(s2).Intersection(s3)

		if badge.Len() != 1 {
			return 0, errors.New("overlap is not one element")
		}

		total += value(badge.ToSlice()[0])
	}

	return total, nil
}

func parseInput(reader io.Reader) ([][]byte, error) {

	var buff bytes.Buffer

	_, err := io.Copy(&buff, reader)
	if err != nil {
		return nil, err
	}

	parts := bytes.Split(buff.Bytes(), []byte{'\n'})

	notEmpty := func(b []byte) bool {
		return len(b) > 0
	}

	filtered := list.Filter(parts, notEmpty)
	return filtered, nil
}

func value(b byte) int {
	if b >= 'a' && b <= 'z' {
		return int(b - 'a' + 1)
	}
	return int(b - 'A' + 27)
}

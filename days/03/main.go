package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/puzzleinput"
	"io"
)

func main() {
	in, err := puzzleinput.Get(3)
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

type rucksack struct {
	content []byte
}

func parseInput(rc io.Reader) ([]rucksack, error) {

	var rucksacks []rucksack

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		rucksacks = append(rucksacks, rucksack{
			content: []byte(line),
		})
	}

	return rucksacks, nil
}

func solve1(rucksacks []rucksack) (int, error) {

	var total int

	for _, r := range rucksacks {
		halve := len(r.content) / 2
		left, right := r.content[:halve], r.content[halve:]

		b := findFirstOverlap(left, right)

		total += value(b)
	}

	return total, nil
}

func value(b byte) int {
	if b >= 'a' && b <= 'z' {
		return int(b - 'a' + 1)
	}
	return int(b - 'A' + 27)
}

func findFirstOverlap(left, right []byte) byte {
	m := make(map[byte]bool)

	for _, b := range left {
		m[b] = true
	}

	for _, b := range right {
		if m[b] {
			return b
		}
	}

	panic("no overlap found")
}

func findAllOverlap(left, right []byte) []byte {
	m := make(map[byte]bool)

	for _, b := range left {
		m[b] = true
	}

	var bs []byte

	for _, b := range right {
		if m[b] {
			bs = append(bs, b)
		}
	}

	return bs
}

func findGroupBadge(rs []rucksack) byte {
	ab := findAllOverlap(rs[0].content, rs[1].content)
	bc := findAllOverlap(rs[1].content, rs[2].content)

	return findFirstOverlap(ab, bc)
}

func solve2(rs []rucksack) (int, error) {
	if len(rs)%3 != 0 {
		panic("not a multiple of 3")
	}

	var total int

	for i := 0; i < len(rs); i += 3 {
		b := findGroupBadge(rs[i : i+3])
		total += value(b)
	}

	return total, nil
}

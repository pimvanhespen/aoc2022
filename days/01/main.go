package main

import (
	"aoc2022/internal/puzzleinput"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input, err := puzzleinput.Get(1)
	if err != nil {
		panic(err)
	}

	invs, err := parseInput(input)
	if err != nil {
		panic(err)
	}

	result1, err := solve1(invs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", result1)

	result2, err := solve2(invs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result2)

}

// solve1 solves the question which Elf has the most calories.
func solve1(elves []Elf) (int, error) {
	var mostCalories int

	for _, elf := range elves {
		calories := elf.Calories()
		if calories > mostCalories {
			mostCalories = calories
		}
	}

	return mostCalories, nil
}

func solve2(elves []Elf) (int, error) {
	sort.Slice(elves, func(i, j int) bool {
		return elves[i].Calories() > elves[j].Calories()
	})

	top3 := elves[0].Calories() + elves[1].Calories() + elves[2].Calories()
	return top3, nil
}

type Elf struct {
	calories int
	items    []int
}

func (i *Elf) Add(item int) {
	i.items = append(i.items, item)
	i.calories += item
}

func (i Elf) Calories() int {
	return i.calories
}

func parseInput(closer io.ReadCloser) (_ []Elf, err error) {
	defer func() {
		cerr := closer.Close()
		if err == nil {
			err = cerr
		}
	}()

	b, err := io.ReadAll(closer)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	var elves []Elf
	var elf Elf

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			elves = append(elves, elf)
			elf = Elf{}
			continue
		}

		n, convErr := strconv.Atoi(line)
		if convErr != nil {
			return nil, fmt.Errorf("bad input on line %d: %w", i, convErr)
		}

		elf.Add(n)
	}

	if len(elf.items) > 0 {
		elves = append(elves, elf)
	}

	return elves, nil
}

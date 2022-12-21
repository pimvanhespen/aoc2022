package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/list"
	"io"
	"strconv"
)

type Item struct {
	Number int
	Index  int
}

func (i Item) String() string {
	return fmt.Sprintf("%d", i.Number)
}

func main() {
	input, err := aoc.Load(20, parse)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))
}

func parse(reader io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(reader)
	var ints []int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		ints = append(ints, num)
	}

	return ints, nil
}

func solve1(input []int) int {
	return solve(input, 1, 1)
}

func solve(input []int, key int, times int) int {
	items := make([]Item, len(input))
	for i, v := range input {
		items[i] = Item{Number: v * key, Index: i}
	}

	c := list.NewLoop[Item](items...)

	mix(c, times)

	return getCoords(c)
}

func solve2(input []int) int {
	const decryptionKey = 811589153
	return solve(input, decryptionKey, 10)
}

func mix(c *list.Loop[Item], times int) {
	maxIndex := c.Size()
	for iter := 0; iter < times; iter++ {
		for i := 0; i < maxIndex; i++ {
			for c.Value().Index != i {
				c.Next()
			}

			item := c.Remove()
			moves := item.Number % (c.Size())

			c.Move(moves)
			c.InsertBefore(item)
		}
	}
}

func getCoords(c *list.Loop[Item]) int {
	for c.Value().Number != 0 {
		c.Next()
	}

	var sum int
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			c.Next()
		}
		sum += c.Value().Number
	}

	return sum
}

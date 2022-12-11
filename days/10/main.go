package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strings"
)

type Instruction struct {
	Op  string
	Arg int
}

func main() {
	r, err := aoc.Get(10)
	if err != nil {
		panic(err)
	}

	in, err := parse(r)
	if err != nil {
		panic(err)
	}

	// Part 1
	s1 := solve1(in)
	fmt.Println("Part 1:", s1)
	//fmt.Println("Part 2:", solve2(in))
}

func parse(reader io.Reader) ([]Instruction, error) {

	var instructions []Instruction

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		if text == "noop" {
			instructions = append(instructions, Instruction{Op: "noop", Arg: 0})
			continue
		}

		var ins Instruction

		_, _ = fmt.Sscanf(text, "%s %d", &ins.Op, &ins.Arg)

		instructions = append(instructions, ins)
	}
	return instructions, nil
}

func solve1(ins []Instruction) int {
	var sum int
	var X int
	var cycle int

	X += 1 // start at 1

	var sb strings.Builder

	eval := func() {
		if abs((cycle%40)-X) <= 1 {
			sb.WriteByte('#')
		} else {
			sb.WriteByte(' ')
		}
		cycle++
		if cycle != 0 && cycle%40 == 0 {
			sb.WriteByte('\n')
		}
		if cycle%40 == 20 {
			sum += (X * cycle)
		}
	}

	for _, in := range ins {
		if in.Op == "noop" {
			eval()
		} else {
			eval()
			eval()
			X += in.Arg
		}
	}

	fmt.Println(sb.String())

	return sum
}

func solve2(in []byte) int {
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/stack"
	"strings"
)

type Harbor struct {
	stacks []stack.Stack[byte]
}

func (h Harbor) Copy() Harbor {
	stacks := make([]stack.Stack[byte], len(h.stacks))
	for n, s := range h.stacks {
		stacks[n] = s.Copy()
	}
	return Harbor{stacks: stacks}
}

func (h Harbor) String() string {
	var b strings.Builder

	var max int
	for _, s := range h.stacks {
		if s.Len() > max {
			max = s.Len()
		}
	}

	var c byte

	for i := 0; i < max; i++ {
		idx := max - i - 1
		_, _ = fmt.Fprintf(&b, "%2d ", idx)

		for _, s := range h.stacks {
			if idx > (s.Len() - 1) {
				b.WriteString("    ")
				continue
			}

			rev := s.Len() - 1 - idx

			c = s.PeekAt(rev)
			_, _ = fmt.Fprintf(&b, "[%c] ", c)
		}

		b.WriteByte('\n')
	}

	b.Write([]byte("   "))

	for i := range h.stacks {
		b.Write([]byte{' ', '1' + byte(i), ' ', ' '})
	}
	b.WriteByte('\n')

	return b.String()
}

type Instruction struct {
	Amount int
	From   int
	To     int
}

func (i Instruction) String() string {
	return fmt.Sprintf("move %d from %d to %d", i.Amount, i.From+1, i.To+1)
}

func main() {

	rc, err := aoc.Get(5)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	scanner := bufio.NewScanner(rc)

	h1, err := ParseHarbor(scanner)
	if err != nil {
		panic(err)
	}

	instructions, err := ParseInstructions(scanner)
	if err != nil {
		panic(err)
	}

	part1 := solve1(h1.Copy(), instructions)
	fmt.Println("Part 1:", part1)

	part2 := solve2(h1.Copy(), instructions)
	fmt.Println("Part 2:", part2)
}

func solve1(harbor Harbor, instructions []Instruction) string {

	var crate byte

	for _, ins := range instructions {
		// Move from top of stack
		for i := 0; i < ins.Amount; i++ {
			crate = harbor.stacks[ins.From].Pop()
			harbor.stacks[ins.To].Push(crate)
		}
	}

	var b strings.Builder

	for _, stack := range harbor.stacks {
		b.WriteByte(stack.Peek())
	}

	return b.String()
}

func solve2(harbor Harbor, instructions []Instruction) string {

	var crates []byte

	for _, ins := range instructions {

		for i := 0; i < ins.Amount; i++ {
			crates = append(crates, harbor.stacks[ins.From].Pop())
		}

		for i := ins.Amount - 1; i >= 0; i-- {
			harbor.stacks[ins.To].Push(crates[i])
		}

		crates = crates[:0] // Clear slice, but preserve capacity (memory)
	}

	var b strings.Builder

	for _, stack := range harbor.stacks {
		b.WriteByte(stack.Peek())
	}

	return b.String()
}

func parseStackRow(line string) []byte {

	var stacks []byte

	for i := 1; i < len(line); i += 4 {
		stacks = append(stacks, line[i])
	}

	return stacks
}

func ParseHarbor(scanner *bufio.Scanner) (Harbor, error) {

	var inputs [][]byte
	var max int

	// Read each horizontal line of stack values
	for scanner.Scan() {
		line := scanner.Text()
		if strings.IndexByte(line, '[') == -1 {
			scanner.Scan()
			break
		}

		row := parseStackRow(line)
		if len(row) > max {
			max = len(row)
		}
		inputs = append(inputs, row)
	}

	// Create stacks
	stacks := make([]stack.Stack[byte], max)
	for i := range stacks {
		stacks[i] = stack.New[byte]()
	}

	// Fill stacks bottom up
	for i := len(inputs) - 1; i >= 0; i-- {
		for j, b := range inputs[i] {
			if b == ' ' {
				continue
			}
			stacks[j].Push(b)
		}
	}

	return Harbor{stacks: stacks}, nil
}

func ParseInstructions(scanner *bufio.Scanner) ([]Instruction, error) {
	var instructions []Instruction

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		var ins Instruction
		_, err := fmt.Sscanf(line, "move %d from %d to %d", &ins.Amount, &ins.From, &ins.To)
		if err != nil {
			return nil, fmt.Errorf("could not parse instruction: %w", err)
		}

		// input isn't zero indexed, the harbor is
		ins.From -= 1
		ins.To -= 1

		instructions = append(instructions, ins)
	}

	return instructions, nil
}

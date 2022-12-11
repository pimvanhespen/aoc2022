package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"sort"
	"strconv"
	"strings"
)

func main() {
	r, err := aoc.Get(11)
	if err != nil {
		panic(err)
	}

	monkeys, err := parse(r)
	if err != nil {
		panic(err)
	}

	part1 := solve1(monkeys, 20)
	fmt.Println("Part 1:", part1)

	for _, monkey := range monkeys {
		monkey.Reset()
	}

	part2 := solve2(monkeys, 10_000, func(i int) int {
		return i / 2
	})
	fmt.Println("Part 2:", part2)
}

type Monkey struct {
	Inspections int
	Op          func(int) int
	Check       func(int) int
	Initial     []int
	Items       []int
}

type Report struct {
	Item   int
	Target int
}

func (m *Monkey) Inspect(adjustWorry func(int) int) (Report, bool) {
	if len(m.Items) == 0 {
		return Report{}, false
	}
	m.Inspections++

	item := m.Items[0]
	m.Items = m.Items[1:]

	item = m.Op(item)
	item = adjustWorry(item)

	return Report{item, m.Check(item)}, true
}

func (m *Monkey) Reset() {
	m.Items = make([]int, len(m.Initial))
	copy(m.Items, m.Initial)
	m.Inspections = 0
}

func (m *Monkey) Receive(item int) {
	m.Items = append(m.Items, item)
}

func (m Monkey) String() string {
	return fmt.Sprintf("Monkey:  inspections=%d, items=%v", m.Inspections, m.Items)
}

func solve1(monkeys []*Monkey, rounds int) int {
	worryFn := func(i int) int { return i / 3 }
	return solve2(monkeys, rounds, worryFn)
}

func solve2(monkeys []*Monkey, rounds int, worryFn func(int) int) int {

	for i := 0; i < rounds; i++ {
		for _, m := range monkeys {
			for {
				report, ok := m.Inspect(worryFn)
				if !ok {
					break
				}
				monkeys[report.Target].Receive(report.Item)
			}
		}
	}

	ints := make([]int, len(monkeys))
	for i, m := range monkeys {
		ints[i] = m.Inspections
	}

	sort.Ints(ints)

	return ints[len(ints)-1] * ints[len(ints)-2]
}

type MathFn func(int, int) int

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func mul(a, b int) int {
	return a * b
}

func div(a, b int) int {
	return a / b
}

func parse(reader io.Reader) ([]*Monkey, error) {

	var monkeys []*Monkey

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		var monkey Monkey

		// first line: items
		scanner.Scan()
		text = scanner.Text()
		listStart := strings.Index(text, ": ")
		strs := strings.Split(text[listStart+2:], ", ")
		for _, str := range strs {
			n, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, err
			}
			monkey.Initial = append(monkey.Initial, int(n))
		}

		// second line: Operation
		scanner.Scan()
		text = scanner.Text()

		var op, arg string
		_, _ = fmt.Sscanf(text, "  Operation: new = old %s %s", &op, &arg)

		var m MathFn

		switch op {
		case "+":
			m = add
		case "-":
			m = sub
		case "*":
			m = mul
		case "/":
			m = div
		}

		n, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			monkey.Op = func(old int) int {
				return m(old, old)
			}
		} else {
			monkey.Op = func(old int) int {
				return m(old, int(n))
			}
		}

		// third line: Check
		scanner.Scan()
		text = scanner.Text()
		var divisor int
		_, _ = fmt.Sscanf(text, "  Test: divisible by %d", &divisor)
		if divisor == 0 {
			return nil, fmt.Errorf("divisor is 0! '%s'", text)
		}
		var whenTrue, whenFalse int
		scanner.Scan()
		text = scanner.Text()
		_, _ = fmt.Sscanf(text, "  If true: throw to monkey %d", &whenTrue)
		scanner.Scan()
		text = scanner.Text()
		_, _ = fmt.Sscanf(text, "  If false: throw to monkey %d", &whenFalse)

		monkey.Check = func(n int) int {
			if n%divisor == 0 {
				return whenTrue
			}
			return whenFalse
		}

		monkey.Reset()
		monkeys = append(monkeys, &monkey)
		scanner.Scan() // skip next line
	}
	return monkeys, nil
}

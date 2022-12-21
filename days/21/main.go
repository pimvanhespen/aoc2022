package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strconv"
	"strings"
)

func main() {
	input, err := aoc.Load(21, parse)
	if err != nil {
		panic(err)
	}

	fmt.Println(solve1(input))
	fmt.Println(solve2(input))
}

type Row struct {
	Name    string
	Left    string
	Right   string
	Op      byte
	Value   int
	IsValue bool
}

func parse(reader io.Reader) ([]Row, error) {
	scanner := bufio.NewScanner(reader)
	var rows []Row
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		parts := strings.Split(text, ": ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid input: %q", text)
		}

		name := parts[0]
		values := strings.Split(parts[1], " ")

		if len(values) == 1 {
			num, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, err
			}
			rows = append(rows, Row{
				Name:    name,
				Value:   num,
				IsValue: true,
			})
			continue
		}

		if len(values) != 3 {
			return nil, fmt.Errorf("invalid input: %q", text)
		}

		left := values[0]
		op := values[1][0]
		right := values[2]

		rows = append(rows, Row{
			Name:  name,
			Left:  left,
			Right: right,
			Op:    op,
		})
	}
	return rows, nil
}

type MathNode struct {
	Name    string
	Left    *MathNode
	Right   *MathNode
	Op      rune
	value   int
	isValue bool
}

func (m MathNode) IsValue() bool {
	return m.isValue
}

func (m MathNode) String() string {
	if m.IsValue() {
		return fmt.Sprintf("%s: %d", m.Name, m.value)
	}
	return fmt.Sprintf("%s: %s %c %s", m.Name, m.Left.Name, m.Op, m.Right.Name)
}

func (m MathNode) Value() int {
	if m.IsValue() {
		return m.value
	}

	switch m.Op {
	case '*':
		return m.Left.Value() * m.Right.Value()
	case '+':
		return m.Left.Value() + m.Right.Value()
	case '-':
		return m.Left.Value() - m.Right.Value()
	case '/':
		return m.Left.Value() / m.Right.Value()
	case '=':
		return m.Left.Value() - m.Right.Value()
	}

	panic(fmt.Sprintf("invalid op: %q", m.Op))
}

func rowsToMap(input []Row) map[string]*MathNode {
	nodes := make(map[string]*MathNode)

	for _, row := range input {

		var node *MathNode
		if n, ok := nodes[row.Name]; ok {
			node = n
		} else {
			node = &MathNode{Name: row.Name}
			nodes[row.Name] = node
		}

		if row.IsValue {
			node.value = int(row.Value)
			node.isValue = true
			continue
		} else {
			node.Op = rune(row.Op)

			var left *MathNode
			if n, ok := nodes[row.Left]; ok {
				left = n
			} else {
				left = &MathNode{Name: row.Left}
				nodes[row.Left] = left
			}

			var right *MathNode
			if n, ok := nodes[row.Right]; ok {
				right = n
			} else {
				right = &MathNode{Name: row.Right}
				nodes[row.Right] = right
			}

			node.Left = left
			node.Right = right
		}
	}

	return nodes
}

func solve1(input []Row) int {

	nodes := rowsToMap(input)

	return nodes["root"].Value()
}

func recursive(node *MathNode, target *MathNode, path []*MathNode) []*MathNode {

	path = append(path, node)

	if node == target {
		return path
	}

	if node.IsValue() {
		return nil
	}

	left := recursive(node.Left, target, path)
	if left != nil {
		return left
	}

	return recursive(node.Right, target, path)
}

func solve2(input []Row) int {
	nodes := rowsToMap(input)

	root, ok := nodes["root"]
	if !ok {
		panic("root not found")
	}
	root.Op = '='

	// find branch leading to human
	human, ok := nodes["humn"]
	if !ok {
		panic("human not found")
	}

	path := recursive(nodes["root"], human, nil)
	if path == nil {
		panic("no path found")
	}

	var toMatch int
	if path[1] == root.Right {
		toMatch = root.Left.Value()
	} else if path[1] == root.Left {
		toMatch = root.Right.Value()
	} else {
		panic("invalid path")
	}

	fmt.Println("to match:", toMatch)
	for idx, node := range path[1 : len(path)-1] {

		var nonPath *MathNode

		isLeft := node.Left == path[idx+2]

		if isLeft {
			nonPath = node.Right
		} else {
			nonPath = node.Left
		}

		oppositeValue := nonPath.Value()

		matcher := reverse(toMatch, node.Op, oppositeValue, isLeft)

		tmp := &MathNode{
			Name:    "tmp",
			value:   matcher,
			isValue: true,
		}

		var old *MathNode

		if isLeft {
			old = node.Left
			node.Left = tmp
		} else {
			old = node.Right
			node.Right = tmp
		}

		fmt.Printf("%s: (M) %15d %c %-15d = %15d ", node.Name, matcher, node.Op, nonPath.Value(), toMatch)
		if v := root.Value(); v >= 1 || v <= -1 {
			fmt.Println("NOK! ", v)
			panic("invalid")
		}
		fmt.Println("(OK)", root.Value())

		if isLeft {
			node.Left = old
		} else {
			node.Right = old
		}

		toMatch = matcher
	}

	human.value = toMatch

	left, right := root.Left.Value(), root.Right.Value()
	if left-right == 0 {
		return int(human.Value())
	}

	panic("no solution found")
}

func reverse(sum int, op rune, value int, isLeft bool) int {
	switch op {
	case '*':
		if isLeft {
			return sum / value
		}
		return sum / value

	case '+':
		return sum - value

	case '-':
		if isLeft {
			return sum + value
		}
		return value - sum

	case '/':
		if isLeft {
			return sum * value
		}
		return value / sum
	}
	panic(fmt.Sprintf("invalid op: %q", op))
}

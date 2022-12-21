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

	// find branch leading to human
	human, ok := nodes["humn"]
	if !ok {
		panic("human not found")
	}

	path := recursive(root, human, nil)
	if path == nil {
		panic("no path found")
	}

	var required int

	for idx, node := range path[:len(path)-1] {

		var opposite int

		isLeft := node.Left == path[idx+1] // +1 for child, +1 for begin offset of path
		if isLeft {
			opposite = node.Right.Value()
		} else {
			opposite = node.Left.Value()
		}

		if idx == 0 {
			// set initial required value to equal whatever the root node has on the opposite side
			required = opposite
			continue
		}

		required = reverse(required, node.Op, opposite, isLeft)
	}

	return required
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

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

type Value interface {
	fmt.Stringer
	IsList() bool
	IntVal() int
	ListVal() ListValue
	Compare(Value) int
}

type IntValue int

func (i IntValue) IsList() bool {
	return false
}

func (i IntValue) IntVal() int {
	return int(i)
}

func (i IntValue) ListVal() ListValue {
	return []Value{i}
}

func (i IntValue) String() string {
	return strconv.Itoa(int(i))
}

func (i IntValue) Compare(other Value) int {
	if other.IsList() {
		return i.ListVal().Compare(other)
	}

	i2 := int(other.IntVal())
	if int(i) == i2 {
		return 0
	}
	if int(i) < i2 {
		return -1
	}
	return 1
}

type ListValue []Value

func (l ListValue) IsList() bool {
	return true
}

func (l ListValue) IntVal() int {
	panic("not an int")
}

func (l ListValue) ListVal() ListValue {
	return l
}

func (l ListValue) String() string {
	var sb strings.Builder

	sb.WriteString("[")

	for i, item := range l {
		sb.WriteString(item.String())
		if i < len(l)-1 {
			sb.WriteString(",")
		}
	}

	sb.WriteString("]")

	return sb.String()
}

func (l ListValue) Compare(other Value) int {
	b := other.ListVal()

	if len(l) == 0 && len(b) == 0 {
		return 0
	}

	if len(l) == 0 {
		return -1
	}

	if len(b) == 0 {
		return 1
	}

	for i := 0; i < len(l) && i < len(b); i++ {
		c := l[i].Compare(b[i])
		if c != 0 {
			return c
		}
	}

	if len(l) == len(b) {
		return 0
	}

	if len(l) < len(b) {
		return -1
	}

	return 1
}

type Packet = ListValue

func main() {
	r, err := aoc.Get(13)
	if err != nil {
		panic(err)
	}

	packets := parse(r)

	// Part 1
	sum1 := solve1(packets)
	println("Part1:", sum1)

	// Part 2
	sum2 := solve2(packets)
	println("Part2:", sum2)
}

func solve1(packets [][2]Packet) int {
	sum := 0

	for n, pair := range packets {
		if pair[0].Compare(pair[1]) < 0 {
			sum += 1 + n
		}
	}

	return sum
}

func solve2(pairs [][2]Packet) int {
	_, divA := parsePacket("[[2]]")
	_, divB := parsePacket("[[6]]")

	packets := make([]Packet, 0, (len(pairs)+1)*2)
	packets = append(packets, divA, divB)

	for _, pair := range pairs {
		packets = append(packets, pair[0], pair[1])
	}

	sort.Slice(packets, func(i, j int) bool {
		return packets[i].Compare(packets[j]) < 0
	})

	total := 1

	for i, p := range packets {
		//fmt.Printf("%-3d: %s\n", i, p)
		if divA.Compare(p) == 0 {
			total *= 1 + i
		} else if divB.Compare(p) == 0 {
			total *= 1 + i
		}
	}

	return total
}

func parse(r io.Reader) [][2]Packet {
	scanner := bufio.NewScanner(r)

	var packets [][2]Packet

	for scanner.Scan() {
		var pair [2]Packet

		for i := range pair {
			s := scanner.Text()
			scanner.Scan() // cycle for next round

			n, p := parsePacket(s)
			if n != len(s) {
				panic(fmt.Sprintf("expected %d, got %d.\n'%s'\n%s", len(s), n, s, p.String()))
			}
			pair[i] = p
		}

		packets = append(packets, pair)
	}

	return packets
}

func parsePacket(s string) (int, Packet) {
	n, p := parseListItem(s[1:])
	return n + 1, p
}

func parseListItem(s string) (int, ListValue) {
	var l ListValue

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '[':
			n, sub := parseListItem(s[i+1:])
			l = append(l, sub)
			i += n
		case ']':
			return i + 1, l
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

			end := strings.IndexAny(s[i:], ",]")
			if end == -1 {
				panic(fmt.Sprintf("no end found for %s", s[i:]))
			}

			n, _ := strconv.Atoi(s[i : i+end])
			l = append(l, IntValue(n))
			i += end - 1
		case ',':
			continue
		default:
		}
	}

	return len(s), l
}

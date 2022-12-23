package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"math"
)

type Vector2D struct {
	X, Y int
}

func (v Vector2D) Add(o Vector2D) Vector2D {
	return Vector2D{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{v.X - other.X, v.Y - other.Y}
}

type Elf struct {
	Pos Vector2D
}

func (e *Elf) Copy() *Elf {
	return NewElf(e.Pos)
}

func (e *Elf) Move(direction Vector2D) {
	e.Pos = e.Pos.Add(direction)
}

func NewElf(pos Vector2D) *Elf {
	return &Elf{
		Pos: pos,
	}
}

type Field struct {
	Elves []*Elf
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) addElf(elf *Elf) {
	f.Elves = append(f.Elves, elf)
}

func (f Field) bounds() (xMin, xMax, yMin, yMax int) {
	xMin, xMax, yMin, yMax = math.MaxInt, math.MinInt, math.MaxInt, math.MinInt

	for _, elf := range f.Elves {
		xMin = min(xMin, elf.Pos.X)
		xMax = max(xMax, elf.Pos.X)
		yMin = min(yMin, elf.Pos.Y)
		yMax = max(yMax, elf.Pos.Y)
	}

	return
}

func (f Field) String() string {
	xMin, xMax, yMin, yMax := f.bounds()

	width, height := xMax-xMin+1, yMax-yMin+1

	bts := make([][]byte, height)
	for i := range bts {
		bts[i] = bytes.Repeat([]byte{'.'}, width)
	}

	for _, elf := range f.Elves {
		bts[elf.Pos.Y-yMin][elf.Pos.X-xMin] = '#'
	}

	return string(bytes.Join(bts, []byte{'\n'}))
}

func (f Field) Copy() *Field {
	field := NewField()
	for _, elf := range f.Elves {
		field.addElf(elf.Copy())
	}
	return field
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parse(reader io.Reader) (*Field, error) {

	bs, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	field := NewField()
	for y, line := range bytes.Split(bs, []byte{'\n'}) {
		for x, c := range line {
			switch c {
			case '#':
				elf := NewElf(Vector2D{X: x, Y: y})
				field.addElf(elf)
			default:
				continue
			}
		}
	}

	return field, nil
}

func (f *Field) Simulate(rounds int, strategy func(field *Field) int) int {
	var total int
	for i := 0; i < rounds; i++ {
		total += strategy(f)
	}
	return total
}

var (
	North     = Vector2D{X: 0, Y: -1}
	NorthEast = Vector2D{X: 1, Y: -1}
	East      = Vector2D{X: 1, Y: 0}
	SouthEast = Vector2D{X: 1, Y: 1}
	South     = Vector2D{X: 0, Y: 1}
	SouthWest = Vector2D{X: -1, Y: 1}
	West      = Vector2D{X: -1, Y: 0}
	NorthWest = Vector2D{X: -1, Y: -1}
)

func Adjecent(d Vector2D) [2]Vector2D {
	if d.X == 0 {
		return [2]Vector2D{
			{X: 1, Y: d.Y},
			{X: -1, Y: d.Y},
		}
	} else {
		return [2]Vector2D{
			{X: d.X, Y: 1},
			{X: d.X, Y: -1},
		}
	}
}

func createStrategy() func(field *Field) int {
	dirs := [4]Vector2D{North, South, West, East}

	return func(field *Field) int {
		elfsMoved := strategy(field, dirs[:])

		dirs[0], dirs[1], dirs[2], dirs[3] = dirs[1], dirs[2], dirs[3], dirs[0]

		return elfsMoved
	}
}

func strategy(field *Field, dirs []Vector2D) int {
	type elfMove struct {
		elf       *Elf
		direction Vector2D
	}

	taken := make(map[Vector2D]struct{})
	for _, elf := range field.Elves {
		taken[elf.Pos] = struct{}{}
	}

	areAnyTaken := func(positions ...Vector2D) bool {
		for _, pos := range positions {
			if _, ok := taken[pos]; ok {
				return true
			}
		}
		return false
	}

	moves := map[Vector2D]elfMove{}
	doubleBooked := map[Vector2D]struct{}{}

	for _, elf := range field.Elves {

		// no neighbours - do nothing
		if !areAnyTaken(
			elf.Pos.Add(North),
			elf.Pos.Add(NorthEast),
			elf.Pos.Add(East),
			elf.Pos.Add(SouthEast),
			elf.Pos.Add(South),
			elf.Pos.Add(SouthWest),
			elf.Pos.Add(West),
			elf.Pos.Add(NorthWest),
		) {
			continue
		}

		for _, direction := range dirs {

			neighbours := Adjecent(direction)
			a, b := neighbours[0], neighbours[1]

			destination := direction.Add(elf.Pos)
			if areAnyTaken(destination, a.Add(elf.Pos), b.Add(elf.Pos)) {
				continue
			}

			if _, reserved := moves[destination]; reserved {
				doubleBooked[destination] = struct{}{}
			} else {
				moves[destination] = elfMove{
					elf:       elf,
					direction: direction,
				}
			}

			break
		}
	}

	for location := range doubleBooked {
		delete(moves, location)
	}

	for _, move := range moves {
		move.elf.Move(move.direction)
	}

	return len(moves)
}

func solve1(field *Field) int {
	const rounds = 10

	field.Simulate(rounds, createStrategy())

	xMin, xMax, yMin, yMax := field.bounds()

	// width * height - totalElvas
	return (xMax-xMin+1)*(yMax-yMin+1) - len(field.Elves)
}

func solve2(field *Field) int {
	strat := createStrategy()

	var rounds int

	for field.Simulate(1, strat) > 0 {
		rounds++
	}
	return 1 + rounds
}

func main() {
	field, err := aoc.Load[*Field](23, parse)
	if err != nil {
		panic(err)
	}

	part1 := solve1(field.Copy())
	fmt.Println("part1:", part1)

	part2 := solve2(field.Copy())
	fmt.Println("part2:", part2)
}

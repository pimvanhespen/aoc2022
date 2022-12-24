package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/set"
	"io"
	"strings"
)

type Valley struct {
	Width     int
	Height    int
	Blizzards []Blizzard
	Occupied  *set.Set[Vector2D]
}

func (v Valley) GetValidMoves(from Vector2D) []Vector2D {
	moves := make([]Vector2D, 0, 5)

	// unrolled loop
	if m := from.Add(North); v.Valid(m) {
		moves = append(moves, m)
	}
	if m := from.Add(South); v.Valid(m) {
		moves = append(moves, m)
	}
	if m := from.Add(East); v.Valid(m) {
		moves = append(moves, m)
	}
	if m := from.Add(West); v.Valid(m) {
		moves = append(moves, m)
	}
	if v.Valid(from) {
		moves = append(moves, from)
	}

	return moves
}

func (v Valley) Valid(next Vector2D) bool {
	if next.X < 0 || next.X >= v.Width {
		return false
	}

	if next.Y < 0 {
		// start
		return next.X == 0 && next.Y == -1
	}

	if next.Y >= v.Height {
		// finish
		return next.X == v.Width-1 && next.Y == v.Height
	}

	return !v.Occupied.Contains(next)
}

func (v Valley) Next() Valley {
	valley := Valley{
		Width:     v.Width,
		Height:    v.Height,
		Blizzards: make([]Blizzard, len(v.Blizzards)),
		Occupied:  set.New[Vector2D](),
	}

	for i, b := range v.Blizzards {
		next := b.Next(v.Width, v.Height)
		valley.Blizzards[i] = next
		valley.Occupied.Add(next.Position)
	}

	return valley
}

func (v Valley) bytes() [][]byte {
	bs := make([][]byte, v.Height+2)
	for i := range bs {
		// regular row
		bs[i] = bytes.Repeat([]byte{Empty}, v.Width+2)
		bs[i][0], bs[i][len(bs[i])-1] = Wall, Wall
	}

	// top and bottom wall
	bs[0] = bytes.Repeat([]byte{Wall}, v.Width+2)
	bs[len(bs)-1] = bytes.Repeat([]byte{Wall}, v.Width+2)

	// start and finish
	bs[0][1] = Empty
	bs[len(bs)-1][len(bs[len(bs)-1])-2] = Empty

	for _, b := range v.Blizzards {
		x, y := 1+b.Position.X, 1+b.Position.Y

		tile := bs[y][x]

		var replacement Tile

		switch tile {
		case Empty:
			replacement = b.Tile()
		case BlizzNorth, BlizzSouth, BlizzEast, BlizzWest:
			replacement = '2'
		case '2', '3', '4', '5', '6', '7', '8':
			replacement = tile + 1
		default:
			replacement = '*'
		}

		bs[y][x] = replacement
	}

	return bs
}

func (v Valley) String() string {
	return string(bytes.Join(v.bytes(), []byte{'\n'}))
}

type Blizzard struct {
	Position  Vector2D
	Direction Vector2D
}

func (b Blizzard) Tile() Tile {
	switch b.Direction {
	case North:
		return BlizzNorth
	case South:
		return BlizzSouth
	case East:
		return BlizzEast
	case West:
		return BlizzWest
	}
	panic("unreachable")
}

func (b Blizzard) Next(width, height int) Blizzard {
	next := b.Position.Add(b.Direction)

	if next.X < 0 {
		next.X = width - 1
	} else if next.X >= width {
		next.X = 0
	} else if next.Y < 0 {
		next.Y = height - 1
	} else if next.Y >= height {
		next.Y = 0
	}

	return Blizzard{
		Position:  next,
		Direction: b.Direction,
	}
}

func NewBlizzard(x, y int, direction Vector2D) Blizzard {
	return Blizzard{
		Position:  NewVector2D(x, y),
		Direction: direction,
	}
}

type Vector2D struct {
	X, Y int
}

func (v Vector2D) Add(o Vector2D) Vector2D {
	return NewVector2D(v.X+o.X, v.Y+o.Y)
}

func NewVector2D(x, y int) Vector2D {
	return Vector2D{
		X: x,
		Y: y,
	}
}

type Tile = byte

const (
	Empty      Tile = '.'
	Wall       Tile = '#'
	Self       Tile = 'E'
	BlizzNorth Tile = '^'
	BlizzSouth Tile = 'v'
	BlizzEast  Tile = '>'
	BlizzWest  Tile = '<'
)

var (
	North = Vector2D{0, -1}
	South = Vector2D{0, 1}
	East  = Vector2D{1, 0}
	West  = Vector2D{-1, 0}
)

func main() {
	valley, err := aoc.Load[Valley](24, parse)
	if err != nil {
		panic(err)
	}

	part1 := solve1(valley)

	fmt.Println("Part 1:", part1)
}

func parse(reader io.Reader) (Valley, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return Valley{}, err
	}

	lines := bytes.Split(b, []byte{'\n'})

	const wallSize = 1

	valley := Valley{
		Width:  len(lines[0]) - 2*wallSize, // remove walls
		Height: len(lines) - 2*wallSize,    // remove walls
	}

	// ignore walls
	for y, line := range lines[1 : len(lines)-1] {
		for x, c := range line[1 : len(line)-1] {
			switch c {
			case BlizzNorth:
				valley.Blizzards = append(valley.Blizzards, NewBlizzard(x, y, North))
			case BlizzSouth:
				valley.Blizzards = append(valley.Blizzards, NewBlizzard(x, y, South))
			case BlizzEast:
				valley.Blizzards = append(valley.Blizzards, NewBlizzard(x, y, East))
			case BlizzWest:
				valley.Blizzards = append(valley.Blizzards, NewBlizzard(x, y, West))
			}
		}
	}

	return valley, nil
}

func solve1(valley Valley) int {

	// initial state
	start := NewVector2D(0, -1)                          // start at the top
	finish := NewVector2D(valley.Width-1, valley.Height) // finish at the bottom

	next := set.New[Vector2D]()
	moves := set.New[Vector2D](start)

	var minutes int

	// for each minute, move all blizzards and check all moves
	// stop when we reach the finish
	for !moves.Contains(finish) {
		if moves.IsEmpty() {
			panic("no path found")
		}

		minutes++
		valley = valley.Next()
		for !moves.IsEmpty() {
			move := moves.Pop()

			for _, m := range valley.GetValidMoves(move) {
				next.Add(m)
			}
		}

		// swap next and moves for the next iteration (saves recreating a set each time)
		moves, next = next, moves
	}

	// find route
	return minutes
}

func FprintState(w io.Writer, valley Valley, moves []Vector2D) {

	b := valley.bytes()

	for _, m := range moves {
		b[m.Y+1][m.X+1] = Self
	}

	var sb strings.Builder

	for _, line := range b {
		for _, c := range line {
			switch c {
			case Empty:
				_, _ = fmt.Fprint(&sb, ".")
			case Wall:
				_, _ = fmt.Fprint(&sb, "#")
				//sb.WriteRune('█')
			case Self:
				_, _ = fmt.Fprint(&sb, "E")
			default:
				_, _ = fmt.Fprintf(&sb, "%c", c)
			}
		}
		_, _ = fmt.Fprintln(&sb)
	}
	_, _ = fmt.Fprintln(&sb)

	_, _ = w.Write([]byte(sb.String()))
}

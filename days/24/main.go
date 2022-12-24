package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/queue"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/set"
	"io"
)

type Valley struct {
	Width      int
	Height     int
	Horizontal *set.Set[State]
	Vertical   *set.Set[State]
}

func (v Valley) IsValid(s State) bool {
	if s.Position.X < 0 || s.Position.X >= v.Width {
		return false
	}

	if s.Position.Y < 0 {
		// is position equal to start position?
		return s.Position == Vector2D{X: 0, Y: -1}
	}

	if s.Position.Y >= v.Height {
		// is position equal to finish?
		return s.Position == Vector2D{X: v.Width - 1, Y: v.Height} // end
	}

	hstate := State{Minute: s.Minute % v.Width, Position: s.Position}
	vstate := State{Minute: s.Minute % v.Height, Position: s.Position}

	return !(v.Vertical.Contains(vstate) || v.Horizontal.Contains(hstate))
}

func (v Valley) bytes(minute int) [][]byte {
	bs := make([][]byte, v.Height+2)
	for i := range bs {
		// regular row
		bs[i] = bytes.Repeat([]byte{'.'}, v.Width+2)
		bs[i][0], bs[i][len(bs[i])-1] = '#', '#'
	}

	// top and bottom wall
	bs[0] = bytes.Repeat([]byte{'#'}, v.Width+2)
	bs[len(bs)-1] = bytes.Repeat([]byte{'#'}, v.Width+2)

	// start and finish
	bs[0][1] = '.'
	bs[len(bs)-1][len(bs[len(bs)-1])-2] = '.'

	for _, b := range v.Horizontal.ToSlice() {
		if b.Minute != minute%v.Width {
			continue
		}

		// + offset
		y, x := 1+b.Position.Y, 1+b.Position.X

		bs[y][x] = '@'
	}

	for _, b := range v.Vertical.ToSlice() {
		if b.Minute != minute%v.Height {
			continue
		}

		y, x := 1+b.Position.Y, 1+b.Position.X

		bs[y][x] = '@'
	}

	return bs
}

func (v Valley) StringWhen(i int) interface{} {
	return string(bytes.Join(v.bytes(i), []byte{'\n'}))
}

type State struct {
	Minute   int
	Position Vector2D
}

func (s State) Equals(o State) bool {
	return s == o
}

func (s State) Options() []State {
	return []State{
		{1 + s.Minute, s.Position}, // wait
		{1 + s.Minute, s.Position.Add(North)},
		{1 + s.Minute, s.Position.Add(East)},
		{1 + s.Minute, s.Position.Add(South)},
		{1 + s.Minute, s.Position.Add(West)},
	}
}

type Vector2D struct {
	X, Y int
}

func NewVector2D(x, y int) Vector2D {
	return Vector2D{
		X: x,
		Y: y,
	}
}

func (v Vector2D) Add(o Vector2D) Vector2D {
	return NewVector2D(v.X+o.X, v.Y+o.Y)
}

func (v Vector2D) Distance(position Vector2D) int {
	return absDiff(v.X, position.X) + absDiff(v.Y, position.Y)
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

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
		Width:      len(lines[0]) - 2*wallSize, // remove walls
		Height:     len(lines) - 2*wallSize,    // remove walls
		Horizontal: set.New[State](),
		Vertical:   set.New[State](),
	}

	calcAllBlizzardStates := func(x, y int, dir Vector2D) []State {
		var states []State

		begin := NewVector2D(x, y)

		max := valley.Width
		if dir == North || dir == South {
			max = valley.Height
		}

		for i := 0; i < max; i++ {
			pos := Vector2D{
				X: (begin.X + i*dir.X + valley.Width) % valley.Width,
				Y: (begin.Y + i*dir.Y + valley.Height) % valley.Height,
			}

			states = append(states, State{
				Minute:   i,
				Position: pos,
			})
		}
		return states
	}

	// ignore walls
	for y, line := range lines[1 : len(lines)-1] {
		for x, c := range line[1 : len(line)-1] {

			switch c {
			case '^':
				valley.Vertical.AddMany(calcAllBlizzardStates(x, y, North)...)
			case 'v':
				valley.Vertical.AddMany(calcAllBlizzardStates(x, y, South)...)
			case '>':
				valley.Horizontal.AddMany(calcAllBlizzardStates(x, y, East)...)
			case '<':
				valley.Horizontal.AddMany(calcAllBlizzardStates(x, y, West)...)
			}
		}
	}

	return valley, nil
}

func solve1(valley Valley) int {

	// initial state
	startTile := NewVector2D(0, -1)                      // start at the top
	finish := NewVector2D(valley.Width-1, valley.Height) // finish at the bottom

	start := State{
		Minute:   0,
		Position: startTile,
	}

	seen := set.New[State]()
	q := queue.NewPriority[State]()
	q.Insert(start, 0)

	// for each minute, move all blizzards and check all moves
	// stop when we reach the finish
	for q.Len() > 0 {

		current := q.Pop()

		if current.Position == finish {
			return current.Minute
		}

		for _, move := range current.Options() {
			// ignore invalid moves, and moves we've already seen
			if (!valley.IsValid(move)) || seen.Contains(move) {
				continue
			}

			seen.Add(move)
			q.Insert(move, current.Minute+finish.Distance(move.Position))
		}
	}

	// find route
	return -1
}

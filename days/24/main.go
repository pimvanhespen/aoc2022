package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/queue"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/set"
	"io"
)

// Valley represents the valley with blizzards
type Valley struct {
	Width      int
	Height     int
	Horizontal *set.Set[State]
	Vertical   *set.Set[State]
}

// IsValid returns true if the given state is valid (e.g. not on a wall / blizzard)
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

// bytes returns the bytes of the valley at the given minute
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

// StringWhen returns the string representation of the valley at the given minute
func (v Valley) StringWhen(i int) interface{} {
	return string(bytes.Join(v.bytes(i), []byte{'\n'}))
}

// State represents a position at a moment in time
type State struct {
	Prev     *State
	Minute   int
	Position Vector2D
}

// Equals returns true if the given state is equal to the current state
func (s State) Equals(o State) bool {
	return s == o
}

// Moves returns all possible moves from the current state
func (s *State) Moves() []State {
	return []State{
		{Prev: s, Minute: 1 + s.Minute, Position: s.Position}, // wait
		{Prev: s, Minute: 1 + s.Minute, Position: s.Position.Add(North)},
		{Prev: s, Minute: 1 + s.Minute, Position: s.Position.Add(East)},
		{Prev: s, Minute: 1 + s.Minute, Position: s.Position.Add(South)},
		{Prev: s, Minute: 1 + s.Minute, Position: s.Position.Add(West)},
	}
}

// Vector2D represents a (x,y)-coordinate
type Vector2D struct {
	X, Y int
}

// NewVector2D returns a new Vector2D
func NewVector2D(x, y int) Vector2D {
	return Vector2D{
		X: x,
		Y: y,
	}
}

// Add returns the sum of the current vector and the given vector
func (v Vector2D) Add(o Vector2D) Vector2D {
	return NewVector2D(v.X+o.X, v.Y+o.Y)
}

// Distance returns the Manhattan distance between the current vector and the given vector
func (v Vector2D) Distance(position Vector2D) int {
	return absDiff(v.X, position.X) + absDiff(v.Y, position.Y)
}

// absDiff the difference between two integers (absolute)
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

	// create a new priority queue
	q := queue.NewPriority[State]()
	q.Insert(start, 0)

	// for each minute, move all blizzards and check all moves
	// stop when we reach the finish
	for q.Len() > 0 {

		// get the most promising state
		current := q.Pop()

		// check if we have reached the end
		if current.Position == finish {
			return current.Minute
		}

		// we are not at the end, so we need to find all possible moves for the current position
		for _, move := range current.Moves() {

			// ignore moves that are not possible, or have been seen before
			if (!valley.IsValid(move)) || seen.Contains(move) {
				continue
			}

			// calculate the priority of the move and insert into the priority queue
			// the moves with the highest priority are the ones that are closest to the finish
			mScore := current.Minute + finish.Distance(move.Position)

			// register the move as seen
			seen.Add(move)
			q.Upsert(move, mScore)
		}
	}

	// find route
	return -1
}

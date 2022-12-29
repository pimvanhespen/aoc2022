package main

import (
	"bytes"
	"reflect"
	"testing"
)

var demo1 = `
   ......
   ......
   ......
   ...
   ...
   ...
......
......
......
...
...
...`

var demo2 = `
     ..........
     ..........
     ..#....#..
     ..........
     ..........
     .....
     .....
     ..#..
     .....
     .....
..........
..........
..#....#..
..........
..........
.....
.....
..#..
.....
.....`

func TestGrid_Move(t *testing.T) {
	size := 5
	grids, err := parseGrid(size, bytes.Split([]byte(demo2), []byte{'\n'})[1:])
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name      string
		grid      *Grid
		position  Vector2D
		direction Vector2D
		moves     []Action
	}

	var trip = []Action{
		Walk(size * 2),
		Right(),
		Walk(1),
		Left(),
		Walk(2),
		Left(),
		Walk(1),
		Right(),
	}

	var roundtrip = append(append(append(trip, trip...), append(trip, trip...)...), Walk(2*size))

	testCases := []testCase{
		{
			name:      "Move Down from Top-Center 1",
			grid:      grids[0],
			position:  Vector2D{X: size / 2, Y: 0},
			direction: Vector2D{X: 0, Y: 1},
			moves:     roundtrip,
		},
		{
			name:      "Move Up from Top-Center 1",
			grid:      grids[0],
			position:  Vector2D{X: size / 2, Y: 0},
			direction: Vector2D{X: 0, Y: -1},
			moves:     roundtrip,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			loc := Location{tc.grid, tc.position, tc.direction}

			moves := DoMoves(loc, tc.moves)

			printGrids(grids, moves)

		})
	}
}
func TestGrid_Move_EmptyBoard(t *testing.T) {
	size := 3
	grids, err := parseGrid(size, bytes.Split([]byte(demo1), []byte{'\n'})[1:])
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name      string
		grid      *Grid
		position  Vector2D
		direction Vector2D
		moves     []Action
	}

	roundtrip := []Action{
		Walk(size),
		Walk(size),
		Walk(size),
		Walk(size),
	}
	center := Vector2D{X: size / 2, Y: size / 2}

	var testCases = []testCase{
		{
			name:      "Move Up from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 0, Y: -1},
			moves:     roundtrip,
		},
		{
			name:      "Move Down from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 0, Y: 1},
			moves:     roundtrip,
		},
		{
			name:      "Move Right from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Left from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: -1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Left from Center 3",
			grid:      grids[2],
			position:  center,
			direction: Vector2D{X: -1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Right from Center 3",
			grid:      grids[2],
			position:  center,
			direction: Vector2D{X: 1, Y: 0},
			moves:     roundtrip,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			loc := Location{tc.grid, tc.position, tc.direction}

			moves := DoMoves(loc, tc.moves)
			printGrids(grids, moves)

		})
	}
}

func TestGrid_Move_EmptyBoard_TestGrid(t *testing.T) {
	size := 4
	grids, err := parseTestGrid(size, bytes.Split([]byte(testEmptyGrid), []byte{'\n'}))
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name      string
		grid      *Grid
		position  Vector2D
		direction Vector2D
		moves     []Action
	}

	roundtrip := []Action{
		Walk(size),
		Walk(size),
		Walk(size),
		Walk(size),
	}
	center := Vector2D{X: size / 2, Y: size / 2}

	var testCases = []testCase{
		{
			name:      "Move Up from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 0, Y: -1},
			moves:     roundtrip,
		},
		{
			name:      "Move Down from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 0, Y: 1},
			moves:     roundtrip,
		},
		{
			name:      "Move Right from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: 1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Left from Center 1",
			grid:      grids[0],
			position:  center,
			direction: Vector2D{X: -1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Left from Center 3",
			grid:      grids[2],
			position:  center,
			direction: Vector2D{X: -1, Y: 0},
			moves:     roundtrip,
		},
		{
			name:      "Move Right from Center 3",
			grid:      grids[2],
			position:  center,
			direction: Vector2D{X: 1, Y: 0},
			moves:     roundtrip,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			loc := Location{tc.grid, tc.position, tc.direction}

			moves := []Location{loc}

			for _, move := range tc.moves {
				newmoves := DoMoves(moves[len(moves)-1], []Action{move})
				printGrids(grids, newmoves)

				moves = append(moves, newmoves...)

			}

			printGrids(grids, moves)

		})
	}
}

func rot90(in [][]byte) [][]byte {
	size := len(in)
	out := make([][]byte, size)
	for i := 0; i < size; i++ {
		out[i] = make([]byte, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			out[x][size-1-y] = in[y][x]
		}
	}

	return out
}

// rotMin90 rotates the grid 90 degrees counter-clockwise
func rotMin90(in [][]byte) [][]byte {
	size := len(in)
	out := make([][]byte, size)
	for i := 0; i < size; i++ {
		out[i] = make([]byte, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			out[size-1-x][y] = in[y][x]
		}
	}

	return out
}

// reverse mirrors the grid horizontally and vertically
func reverse(in [][]byte) [][]byte {
	size := len(in)
	out := make([][]byte, size)
	for i := 0; i < size; i++ {
		out[i] = make([]byte, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			out[size-1-y][size-1-x] = in[y][x]
		}
	}

	return out
}

func TestReverse(t *testing.T) {
	in := [][]byte{
		{1, 2},
		{3, 4},
	}
	expect := [][]byte{
		{4, 3},
		{2, 1},
	}

	out := reverse(in)
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("Expected %v, got %v", expect, out)
	}
}

func TestRot90(t *testing.T) {
	in := [][]byte{
		{1, 2},
		{3, 4},
	}
	expect := [][]byte{
		{3, 1},
		{4, 2},
	}

	out := rot90(in)
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("rot90(%v) = %v, expected %v", in, out, expect)
	}
}

func TestRotMin90(t *testing.T) {
	in := [][]byte{
		{1, 2},
		{3, 4},
	}
	expect := [][]byte{
		{2, 4},
		{1, 3},
	}

	out := rotMin90(in)
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("got %v, expected %v", out, expect)
	}
}

const testEmptyGrid = `        ....
        ....
        ....
        ....
............
............
............
............
        ........
        ........
        ........
        ........`

func parseTestGrid(size int, data [][]byte) ([]*Grid, error) {
	makeGrid := func(pos Vector2D) *Grid {
		offset := pos.Factor(size)

		grid := &Grid{
			Data:   make([][]byte, size),
			Size:   size,
			Offset: offset,
		}
		for y := range grid.Data {
			grid.Data[y] = data[offset.Y+y][offset.X : offset.X+size]
		}
		return grid
	}

	var (
		one   = makeGrid(Vector2D{2, 0})
		two   = makeGrid(Vector2D{0, 1})
		three = makeGrid(Vector2D{1, 1})
		four  = makeGrid(Vector2D{2, 1})
		five  = makeGrid(Vector2D{2, 2})
		six   = makeGrid(Vector2D{3, 2})
	)

	linkTestGrids(one, two, three, four, five, six)

	return []*Grid{one, two, three, four, five, six}, nil
}

func linkTestGrids(one, two, three, four, five, six *Grid) {
	sizeMask := len(one.Data) - 1

	invert := func(v int) int {
		return sizeMask - v
	}

	one.Top = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.X), Y: 0}
		},
		Direction: DirDown,
	}

	one.Bottom = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: v.X, Y: 0}
		},
		Direction: DirDown,
	}

	one.Left = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: v.Y, Y: 0}
		},
		Direction: DirDown,
	}

	one.Right = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: invert(v.Y)}
		},
		Direction: DirLeft,
	}

	two.Top = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.X), Y: 0}
		},
		Direction: DirDown,
	}

	two.Bottom = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.X), Y: sizeMask}
		},
		Direction: DirUp,
	}

	two.Left = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.Y), Y: sizeMask}
		},
		Direction: DirUp,
	}

	two.Right = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: v.Y}
		},
		Direction: DirRight,
	}

	three.Top = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: v.X}
		},
		Direction: DirRight,
	}

	three.Bottom = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: invert(v.X)}
		},
		Direction: DirRight,
	}

	three.Left = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: v.Y}
		},
		Direction: DirLeft,
	}

	three.Right = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: v.Y}
		},
		Direction: DirRight,
	}

	four.Top = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: v.X, Y: sizeMask}
		},
		Direction: DirUp,
	}

	four.Bottom = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: v.X, Y: 0}
		},
		Direction: DirDown,
	}

	four.Left = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: v.Y}
		},
		Direction: DirLeft,
	}

	four.Right = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.Y), Y: 0}
		},
		Direction: DirDown,
	}

	five.Top = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: v.X, Y: sizeMask}
		},
		Direction: DirUp,
	}

	five.Bottom = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.X), Y: sizeMask}
		},
		Direction: DirUp,
	}

	five.Left = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: invert(v.Y), Y: sizeMask}
		},
		Direction: DirUp,
	}

	five.Right = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: v.Y}
		},
		Direction: DirRight,
	}

	six.Top = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: invert(v.X)}
		},
		Direction: DirLeft,
	}

	six.Bottom = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: 0, Y: invert(v.X)}
		},
		Direction: DirRight,
	}

	six.Left = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: v.Y}
		},
		Direction: DirLeft,
	}

	six.Right = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{X: sizeMask, Y: invert(v.Y)}
		},
		Direction: DirLeft,
	}
}

//func TestDemoInput(t *testing.T) {
//
//	const size = 4
//
//	data := bytes.Split([]byte(testData), []byte("\n"))
//
//	makeGrid := func(x, y int) *Grid {
//		grid := &Grid{
//			Data:   make([][]byte, size),
//			Size:   size,
//			Offset: Vector2D{x, y},
//		}
//
//		for n := range grid.Data {
//			grid.Data[n] = data[y+n][x : x+size]
//		}
//
//		return grid
//	}
//
//	a := makeGrid(8, 0)
//	b := makeGrid(0, 0)
//
//	grids, err := parseGrid(4, bytes.Split([]byte(testEmptyGrid), []byte{'\n'}))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	steps, err := parseSteps([]byte("10R5L5R10L4R5L5"))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	start := Location{
//		grid: grids[0],
//		pos:  Vector2D{X: 0, Y: 0},
//		dir:  DirRight,
//	}
//
//	moves := DoMoves(start, steps)
//
//	printGrids(grids, moves)
//
//	last := moves[len(moves)-1]
//	fmt.Println("Final Location:", last)
//
//	var result int
//	result += 1000 * (1 + last.pos.Y + last.grid.Offset.Y)
//	result += 4 * (1 + last.pos.X + last.grid.Offset.X)
//	switch last.dir {
//	case DirDown:
//		result += 1
//	case DirLeft:
//		result += 2
//	case DirUp:
//		result += 3
//		//case DirRight:
//		//	result += 0
//	}
//
//	fmt.Println(result)
//}

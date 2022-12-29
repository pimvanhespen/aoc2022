package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	grids, route, err := parse(strings.NewReader(strings.Trim(input, "\n")))
	if err != nil {
		panic(err)
	}

	direction := Vector2D{X: 1, Y: 0}
	position := Vector2D{X: 0, Y: 0}
	current := grids[0]

	start := Location{grid: current, pos: position, dir: direction}

	moves := DoMoves(start, route[:])

	printGrids(grids, moves)

	last := moves[len(moves)-1]

	fmt.Println("Last location", last.pos)

	var result int
	result += 1000 * (1 + last.pos.Y + last.grid.Offset.Y)
	result += 4 * (1 + last.pos.X + last.grid.Offset.X)
	switch last.dir {
	case DirDown:
		result += 1
	case DirLeft:
		result += 2
	case DirUp:
		result += 3
		//case DirRight:
		//	result += 0
	}

	fmt.Println(result)
}

type Location struct {
	grid *Grid
	pos  Vector2D
	dir  Vector2D
}

func (l Location) Do(act ActionType) (Location, bool) {
	switch act {
	case RotateLeft:
		l.dir = rotateLeft(l.dir)
	case RotateRight:
		l.dir = rotateRight(l.dir)
	case WalkDistance:
		ng, np, nd, ok := l.grid.Move(l.pos, l.dir)
		if !ok {
			return l, false // no move
		}
		return Location{grid: ng, pos: np, dir: nd}, true
	}
	return l, true
}

type Vector2D struct {
	X, Y int
}

func (v Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2D) Factor(size int) Vector2D {
	return Vector2D{
		X: v.X * size,
		Y: v.Y * size,
	}
}

type Grid struct {
	Data   [][]byte
	Size   int
	Offset Vector2D
	Top    *Transition
	Bottom *Transition
	Left   *Transition
	Right  *Transition
}

func (g *Grid) Move(position Vector2D, direction Vector2D) (*Grid, Vector2D, Vector2D, bool) {
	newPosition := position.Add(direction)

	var transition *Transition

	if newPosition.Y < 0 {
		transition = g.Top
	} else if newPosition.Y >= g.Size {
		transition = g.Bottom
	} else if newPosition.X < 0 {
		transition = g.Left
	} else if newPosition.X >= g.Size {
		transition = g.Right
	} else {
		// stay on current grid
		if g.IsWall(newPosition) {
			return g, position, direction, false
		}
		return g, newPosition, direction, true
	}

	nextPos := transition.Translate(position)

	if transition.Grid.IsWall(nextPos) {
		return g, position, direction, false
	}

	return transition.Grid, transition.Translate(position), transition.Direction, true
}

func (g *Grid) IsWall(pos Vector2D) bool {
	return g.Data[pos.Y][pos.X] == '#'
}

type Transition struct {
	Grid      *Grid
	Translate func(Vector2D) Vector2D
	Direction Vector2D
}

var (
	DirLeft  = Vector2D{X: -1}
	DirRight = Vector2D{X: 1}
	DirUp    = Vector2D{Y: -1}
	DirDown  = Vector2D{Y: 1}
)

func rotateLeft(v Vector2D) Vector2D {
	switch v {
	case DirLeft:
		return DirDown
	case DirRight:
		return DirUp
	case DirUp:
		return DirLeft
	case DirDown:
		return DirRight
	}
	panic("invalid direction")
}

func rotateRight(v Vector2D) Vector2D {
	switch v {
	case DirLeft:
		return DirUp
	case DirRight:
		return DirDown
	case DirUp:
		return DirRight
	case DirDown:
		return DirLeft
	}
	panic("invalid direction")
}

type ActionType int

const (
	Unknown ActionType = iota
	WalkDistance
	RotateLeft
	RotateRight
)

type Action struct {
	distance int
	Type     ActionType
}

func (a Action) String() string {
	switch a.Type {
	case WalkDistance:
		return strconv.Itoa(a.distance)
	case RotateLeft:
		return "L"
	case RotateRight:
		return "R"
	}
	return "?"
}

func Walk(distance int) Action {
	return Action{distance: distance, Type: WalkDistance}
}

func Left() Action {
	return Action{Type: RotateLeft}
}

func Right() Action {
	return Action{Type: RotateRight}
}

func parse(reader io.Reader) ([]*Grid, []Action, error) {

	bts, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	lines := bytes.Split(bts, []byte{'\n'})

	gridData := lines[0 : len(lines)-2]

	fmt.Println(string(bytes.Join(gridData, []byte{'\n'})))

	size := len(lines) / 4

	grids, err := parseGrid(size, gridData)
	if err != nil {
		return nil, nil, err
	}

	steps, err := parseSteps(lines[len(lines)-1])
	if err != nil {
		return nil, nil, err
	}

	return grids, steps, nil
}

func parseSteps(bts []byte) ([]Action, error) {
	var steps []Action

	input := string(bts)

	for len(bts) > 0 {
		switch bts[0] {
		case 'L':
			steps = append(steps, Left())
			bts = bts[1:]
		case 'R':
			steps = append(steps, Right())
			bts = bts[1:]
		default:
			nextLR := bytes.IndexAny(bts, "LR")
			if nextLR == -1 {
				nextLR = len(bts)
			}

			n, err := strconv.Atoi(string(bts[:nextLR]))
			if err != nil {
				return nil, err
			}
			steps = append(steps, Walk(n))
			bts = bts[nextLR:]
		}
	}

	var sb strings.Builder
	for _, step := range steps {
		sb.WriteString(step.String())
	}

	if sb.String() != input {
		panic("invalid parse")
	}

	return steps, nil
}

func linkGrids(one, two, three, four, five, six *Grid) {
	sizeMask := one.Size - 1

	invert := func(v int) int {
		return sizeMask - v
	}

	// 1 -> Right
	one.Right = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{0, v.Y}
		},
		Direction: Vector2D{1, 0},
	}
	// 1 -> Left
	one.Left = &Transition{
		Grid: five,
		Translate: func(d Vector2D) Vector2D {
			return Vector2D{d.X, invert(d.Y)}
		},
		Direction: Vector2D{1, 0},
	}
	// 1 -> Bottom
	one.Bottom = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 1 -> Top
	one.Top = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{0, v.X}
		},
		Direction: Vector2D{1, 0},
	}

	// 2 -> Right
	two.Right = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, invert(v.Y)}
		},
		Direction: Vector2D{-1, 0},
	}
	// 2 -> Left
	two.Left = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, v.Y}
		},
		Direction: Vector2D{-1, 0},
	}
	// 2 -> Bottom
	two.Bottom = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, v.X}
		},
		Direction: Vector2D{-1, 0},
	}
	// 2 -> Top
	two.Top = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}

	// 3 -> Right
	three.Right = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.Y, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}
	// 3 -> Left
	three.Left = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.Y, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 3 -> Bottom
	three.Bottom = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 3 -> Top
	three.Top = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}

	// 4 -> Right
	four.Right = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, invert(v.Y)}
		},
		Direction: Vector2D{-1, 0},
	}
	// 4 -> Left
	four.Left = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, v.Y}
		},
		Direction: Vector2D{-1, 0},
	}
	// 4 -> Bottom
	four.Bottom = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{sizeMask, v.X}
		},
		Direction: Vector2D{-1, 0},
	}
	// 4 -> Top
	four.Top = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}

	// 5 -> Right
	five.Right = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{0, v.Y}
		},
		Direction: Vector2D{1, 0},
	}
	// 5 -> Left
	five.Left = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{0, invert(v.Y)}
		},
		Direction: Vector2D{1, 0},
	}
	// 5 -> Bottom
	five.Bottom = &Transition{
		Grid: six,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 5 -> Top
	five.Top = &Transition{
		Grid: three,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{0, v.X}
		},
		Direction: Vector2D{1, 0},
	}

	// 6 -> Right
	six.Right = &Transition{
		Grid: four,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.Y, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}
	// 6 -> Left
	six.Left = &Transition{
		Grid: one,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.Y, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 6 -> Bottom
	six.Bottom = &Transition{
		Grid: two,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, 0}
		},
		Direction: Vector2D{0, 1},
	}
	// 6 -> Top
	six.Top = &Transition{
		Grid: five,
		Translate: func(v Vector2D) Vector2D {
			return Vector2D{v.X, sizeMask}
		},
		Direction: Vector2D{0, -1},
	}
}

func parseGrid(size int, data [][]byte) ([]*Grid, error) {

	// shape
	//   XX
	//   X
	//  XX
	//  X

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
		one   = makeGrid(Vector2D{1, 0})
		two   = makeGrid(Vector2D{2, 0})
		three = makeGrid(Vector2D{1, 1})
		four  = makeGrid(Vector2D{1, 2})
		five  = makeGrid(Vector2D{0, 2})
		six   = makeGrid(Vector2D{0, 3})
	)

	linkGrids(one, two, three, four, five, six)

	// stitch together

	return []*Grid{one, two, three, five, four, six}, nil
}

func DoMoves(loc Location, moves []Action) []Location {
	var locs []Location

	locs = append(locs, loc)

	for n, action := range moves {
		log.Println("Move", n, action)

		switch action.Type {
		case RotateLeft:
			loc, _ = loc.Do(RotateLeft)
		case RotateRight:
			loc, _ = loc.Do(RotateRight)
		case WalkDistance:
			for i := 0; i < action.distance; i++ {
				next, ok := loc.Do(WalkDistance)
				if !ok {
					break
				}
				loc = next
				locs = append(locs, loc)
			}
		}
	}

	return locs
}

//func printGrids(grids []*Grid, moves []Location) {
//	size := len(grids[0].Data[0])
//	var sb strings.Builder
//
//	gridBytes := [6][][]byte{}
//
//	for i, grid := range grids {
//		gridBytes[i] = make([][]byte, size)
//		for j := range gridBytes[i] {
//			gridBytes[i][j] = make([]byte, len(grid.Data[j]))
//			copy(gridBytes[i][j], grid.Data[j])
//		}
//
//		for _, m := range moves {
//			if m.grid == grid {
//
//				var c byte
//
//				//c = '0' + byte(n%10)
//
//				switch m.dir {
//				case Vector2D{X: 1, Y: 0}:
//					c = '>'
//				case Vector2D{X: -1, Y: 0}:
//					c = '<'
//				case Vector2D{X: 0, Y: 1}:
//					c = 'v'
//				case Vector2D{X: 0, Y: -1}:
//					c = '^'
//				default:
//					panic("unknown direction")
//				}
//
//				gridBytes[i][m.pos.Y][m.pos.X] = c
//			}
//		}
//	}
//
//	for i := 0; i < size; i++ {
//		sb.WriteString(strings.Repeat(" ", size))
//		sb.WriteString(strings.Repeat(" ", size))
//		sb.WriteString(string(gridBytes[0][i]))
//		sb.WriteByte('\n')
//	}
//
//	for i := 0; i < size; i++ {
//		sb.WriteString(string(gridBytes[1][i]))
//		sb.WriteString(string(gridBytes[2][i]))
//		sb.WriteString(string(gridBytes[3][i]))
//		sb.WriteByte('\n')
//	}
//
//	for i := 0; i < size; i++ {
//		sb.WriteString(strings.Repeat(" ", size))
//		sb.WriteString(strings.Repeat(" ", size))
//		sb.WriteString(string(gridBytes[4][i]))
//		sb.WriteString(string(gridBytes[5][i]))
//		sb.WriteByte('\n')
//	}
//
//	println(sb.String())
//}

func printGrids(grids []*Grid, moves []Location) {
	size := len(grids[0].Data[0])
	var sb strings.Builder

	gridBytes := [6][][]byte{}

	for i, grid := range grids {
		gridBytes[i] = make([][]byte, size)
		for j := range gridBytes[i] {
			gridBytes[i][j] = make([]byte, len(grid.Data[j]))
			copy(gridBytes[i][j], grid.Data[j])
		}

		for _, m := range moves {
			if m.grid == grid {

				var c byte

				//c = '0' + byte(n%10)

				switch m.dir {
				case Vector2D{X: 1, Y: 0}:
					c = '>'
				case Vector2D{X: -1, Y: 0}:
					c = '<'
				case Vector2D{X: 0, Y: 1}:
					c = 'v'
				case Vector2D{X: 0, Y: -1}:
					c = '^'
				default:
					panic("unknown direction")
				}

				gridBytes[i][m.pos.Y][m.pos.X] = c
			}
		}
	}

	for i := 0; i < size; i++ {
		sb.WriteString(strings.Repeat(" ", size))
		sb.WriteString(string(gridBytes[0][i]))
		sb.WriteString(string(gridBytes[1][i]))
		sb.WriteByte('\n')
	}

	for i := 0; i < size; i++ {
		sb.WriteString(strings.Repeat(" ", size))
		sb.WriteString(string(gridBytes[2][i]))
		sb.WriteByte('\n')
	}

	for i := 0; i < size; i++ {
		sb.WriteString(string(gridBytes[3][i]))
		sb.WriteString(string(gridBytes[4][i]))
		sb.WriteByte('\n')
	}

	for i := 0; i < size; i++ {
		sb.WriteString(string(gridBytes[5][i]))
		sb.WriteByte('\n')
	}

	println(sb.String())
}

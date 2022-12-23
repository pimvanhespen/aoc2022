package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strconv"
	"strings"
)

func main() {
	r, err := aoc.Get(22)
	if err != nil {
		panic(err)
	}

	field, steps, err := parse(r)
	if err != nil {
		panic(err)
	}

	r1 := solve1(field.Copy(), steps)
	fmt.Println("Part 1:", r1)
}

func solve1(field *Field, steps []Step) int {

	begin := field.start
	direction := Right

	err := field.JoinEdges(joinSimple)
	if err != nil {
		panic(err)
	}

	place, history := executeRoute(begin, direction, steps)

	fmt.Println(plotMoves(field, history))

	return calcScore1(place, history)
}

func parseField(data [][]byte) (*Field, error) {
	field := NewField()

	for y, line := range data {
		for x, c := range line {

			var t TileType
			pos := Vector2D{x, y}

			switch c {
			case ' ':
				continue
			default:
				t = TileType(c)
			}

			tile := NewTile(t, pos)

			// initial connections
			if left, ok := field.tiles[pos.Add(Left)]; ok {
				tile.Left = left
				left.Right = tile
			}

			if top, ok := field.tiles[pos.Add(Up)]; ok {
				tile.Top = top
				top.Bottom = tile
			}

			field.register(tile)
		}
	}
	return field, nil
}

func parse(reader io.Reader) (*Field, []Step, error) {
	var data [][]byte

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		data = append(data, []byte(line))
	}

	field, err := parseField(data)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse field: %w", err)
	}

	scanner.Scan()
	route := scanner.Text()
	if route == "" {
		return nil, nil, fmt.Errorf("no route found")
	}

	steps, err := parseSteps([]byte(route))
	if err != nil {
		return nil, nil, err
	}

	return field, steps, nil
}

func parseSteps(bts []byte) ([]Step, error) {
	var turns int
	turns += bytes.Count(bts, []byte{'L'})
	turns += bytes.Count(bts, []byte{'R'})

	steps := make([]Step, 0, turns*2+1) // each turn is paired with forwards steps, maybe the end isn't

	for len(bts) > 0 {

		switch bts[0] {
		case 'R':
			steps = append(steps, Step{rotateClockwise: true})
			bts = bts[1:]
			continue
		case 'L':
			steps = append(steps, Step{rotateClockwise: false})
			bts = bts[1:]
			continue
		default:
		}

		ri, li := bytes.IndexByte(bts, 'R'), bytes.IndexByte(bts, 'L')

		var offset int
		if ri == -1 && li == -1 {
			offset = len(bts)
		} else if ri == -1 {
			offset = li
		} else if li == -1 {
			offset = ri
		} else {
			offset = min(ri, li)
		}

		distance, err := strconv.ParseInt(string(bts[:offset]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse stepsize `%s`: %w", string(bts), err)
		}
		bts = bts[offset:]

		steps = append(steps, Step{
			isMove:   true,
			distance: int(distance),
		})
	}

	return steps, nil
}

func executeRoute(place *Tile, direction Vector2D, steps []Step) (*Tile, map[Vector2D]Vector2D) {

	history := make(map[Vector2D]Vector2D)
	history[place.Pos] = direction

	for _, s := range steps {
		if s.isMove {
			for i := 0; i < s.distance; i++ {
				next := place.Move(direction)
				if next.Type == Wall {
					break
				}
				place = next
				history[place.Pos] = direction
			}
		} else {
			direction = Rotate(direction, s.rotateClockwise)
		}
		history[place.Pos] = direction
	}

	return place, history
}

func joinSimple(field *Field) error {

	w, h := field.Dimensions()

	tops := make([]*Tile, w)
	bottoms := make([]*Tile, w)

	lefts := make([]*Tile, h)
	rights := make([]*Tile, h)

	// stitch all outer tiles to their opposite neighbours (wrap around the field)
	for x := field.xMin; x <= field.xMax; x++ {
		for y := field.yMin; y <= field.yMax; y++ {
			pos := Vector2D{x, y}

			tile, ok := field.tiles[pos]
			if !ok {
				continue
			}
			if tile.Top == nil {
				if tops[x] != nil {
					return fmt.Errorf("top tile already set for x=%d", x)
				}
				tops[x] = tile
			}
			if tile.Bottom == nil {
				if bottoms[x] != nil {
					return fmt.Errorf("bottom tile already set for x=%d", x)
				}
				bottoms[x] = tile
			}
			if tile.Left == nil {
				if lefts[y] != nil {
					return fmt.Errorf("left tile already set for y=%d", y)
				}
				lefts[y] = tile
			}
			if tile.Right == nil {
				if rights[y] != nil {
					return fmt.Errorf("right tile already set for y=%d", y)
				}
				rights[y] = tile
			}
		}
	}

	if len(tops) != len(bottoms) {
		return fmt.Errorf("tops and bottoms do not match")
	}
	if len(lefts) != len(rights) {
		return fmt.Errorf("lefts and rights do not match")
	}

	for i := 0; i < len(tops); i++ {
		tops[i].Top = bottoms[i]
		bottoms[i].Bottom = tops[i]
	}

	for i := 0; i < len(lefts); i++ {
		lefts[i].Left = rights[i]
		rights[i].Right = lefts[i]
	}

	return nil
}

func calcScore1(t *Tile, history map[Vector2D]Vector2D) int {

	var score int

	score += (1 + t.Pos.Y) * 1_000
	score += (1 + t.Pos.X) * 4

	direction, ok := history[t.Pos]
	if !ok {
		panic("no direction for last step")
	}

	switch direction {
	case Right:
		score += 0
	case Down:
		score += 1
	case Left:
		score += 2
	case Up:
		score += 3
	}
	return score
}

type TileType byte

const (
	Invalid TileType = 0
	Floor            = '.'
	Wall             = '#'
)

type Tile struct {
	Top    *Tile
	Bottom *Tile
	Left   *Tile
	Right  *Tile
	Pos    Vector2D
	Type   TileType
}

func (t *Tile) String() string {
	return fmt.Sprintf("%c", t.Type)
}

func NewTile(t TileType, pos Vector2D) *Tile {
	if t == Invalid {
		panic("invalid tile type")
	}
	return &Tile{
		Top:    nil,
		Bottom: nil,
		Left:   nil,
		Right:  nil,
		Pos:    pos,
		Type:   t,
	}
}

func (t Tile) Move(direction Vector2D) *Tile {
	switch direction {
	case Up:
		return t.Top
	case Down:
		return t.Bottom
	case Left:
		return t.Left
	case Right:
		return t.Right
	}
	panic("invalid direction")
}

func (t Tile) Copy() *Tile {
	return NewTile(t.Type, t.Pos)
}

type Field struct {
	start                  *Tile
	xMin, xMax, yMin, yMax int
	tiles                  map[Vector2D]*Tile
}

func NewField() *Field {
	return &Field{tiles: make(map[Vector2D]*Tile)}
}

func (f *Field) register(t *Tile) {
	if f.start == nil {
		f.start = t
	}
	f.xMin = min(f.xMin, t.Pos.X)
	f.xMax = max(f.xMax, t.Pos.X)
	f.yMin = min(f.yMin, t.Pos.Y)
	f.yMax = max(f.yMax, t.Pos.Y)
	f.tiles[t.Pos] = t
}

func (f Field) Dimensions() (int, int) {
	return f.xMax - f.xMin + 1, f.yMax - f.yMin + 1
}

func (f *Field) JoinEdges(strategy func(field *Field) error) error {
	return strategy(f)
}

func (f Field) Copy() *Field {
	nf := NewField()

	for _, tile := range f.tiles {
		c := tile.Copy()
		nf.register(c)
	}

	nf.start = nf.tiles[f.start.Pos]

	for _, tile := range nf.tiles {
		// get left and top neighbours
		left, ok := nf.tiles[tile.Pos.Add(Left)]
		if ok {
			tile.Left = left
			left.Right = tile
		}

		top, ok := nf.tiles[tile.Pos.Add(Up)]
		if ok {
			tile.Top = top
			top.Bottom = tile
		}
	}

	return nf
}

func (f *Field) Get(d Vector2D) (*Tile, bool) {
	t, ok := f.tiles[d]
	return t, ok
}

var (
	Up    = Vector2D{0, -1}
	Down  = Vector2D{0, 1}
	Left  = Vector2D{-1, 0}
	Right = Vector2D{1, 0}
)

type Vector2D struct {
	X, Y int
}

func (v Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v.X + v2.X, v.Y + v2.Y}
}

func Rotate(v Vector2D, clockwise bool) Vector2D {
	switch v {
	case Up:
		if clockwise {
			return Right
		}
		return Left
	case Right:
		if clockwise {
			return Down
		}
		return Up
	case Down:
		if clockwise {
			return Left
		}
		return Right
	case Left:
		if clockwise {
			return Up
		}
		return Down
	}
	panic("invalid direction")
}

type Step struct {
	isMove          bool
	rotateClockwise bool
	distance        int
}

func plotMoves(field *Field, history map[Vector2D]Vector2D) string {
	w, h := field.Dimensions()

	var b strings.Builder

	green := color.New(color.FgGreen)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := Vector2D{x, y}
			if dir, HaveBeenHere := history[v]; HaveBeenHere {
				var d byte
				switch dir {
				case Up:
					d = '^'
				case Down:
					d = 'v'
				case Left:
					d = '<'
				case Right:
					d = '>'
				}
				_, _ = green.Fprintf(&b, "%c", d)
			} else if t, IsPartOfField := field.tiles[v]; IsPartOfField {
				b.WriteByte(byte(t.Type))
			} else {
				b.WriteByte(' ')
			}
		}
		b.WriteByte('\n')
	}

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

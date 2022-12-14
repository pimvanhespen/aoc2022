package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/list"
	"io"
	"math"
	"strings"
)

type FieldBuilder struct {
	xMin, xMax, yMin, yMax int
	Lines                  []Line
}

type Field struct {
	xMin  int
	xMax  int
	yMin  int
	yMax  int
	field [][]byte
}

func NewField() FieldBuilder {
	return FieldBuilder{
		xMin: math.MaxInt,
		xMax: 0,
		yMin: math.MaxInt,
		yMax: 0,
	}
}

func (f *FieldBuilder) Add(lines []Line) {
	for _, l := range lines {
		f.xMin = min(f.xMin, min(l.From.X, l.To.X))
		f.xMax = max(f.xMax, max(l.From.X, l.To.X))
		f.yMin = min(f.yMin, min(l.From.Y, l.To.Y))
		f.yMax = max(f.yMax, max(l.From.Y, l.To.Y))

		f.Lines = append(f.Lines, l)
	}
}

func (f FieldBuilder) Compile() Field {
	bts := make([][]byte, f.yMax+1)
	for i := range bts {
		bts[i] = bytes.Repeat([]byte{'.'}, f.xMax-f.xMin+1)
	}

	for _, l := range f.Lines {
		for x := l.From.X; x <= l.To.X; x++ {
			for y := l.From.Y; y <= l.To.Y; y++ {
				xo, yo := x-f.xMin, y
				bts[yo][xo] = '#'
			}
		}
	}

	return Field{
		xMin:  f.xMin,
		xMax:  f.xMax,
		yMin:  f.yMin,
		yMax:  f.yMax,
		field: bts,
	}
}

func (f Field) String() string {

	width, height := f.xMax-f.xMin+1, f.yMax+1

	bts := make([][]byte, height)
	for i := range bts {
		bts[i] = append(bytes.Repeat([]byte{' '}, 4), f.field[i]...)
		if i%10 == 0 {
			s := fmt.Sprintf("%3d", i)
			bts[i][0] = s[0]
			bts[i][1] = s[1]
			bts[i][2] = s[2]
		}
	}

	bts[0][500-f.xMin+4] = 'o'

	header := make([][]byte, 4)
	for i := range header {
		header[i] = bytes.Repeat([]byte{' '}, width+4)
	}

	for i := 4; i < width+4; i++ {
		if (i-4+f.xMin)%10 == 0 || i == 4 || i == width+3 {
			n := fmt.Sprintf("%3d", i-4+f.xMin)
			header[0][i] = n[0]
			header[1][i] = n[1]
			header[2][i] = n[2]
		}
	}

	bts = append(header, bts...)

	return string(bytes.Join(bts, []byte{'\n'}))
}

func (f Field) isAir(pos Vector) bool {
	if pos.X < 0 || pos.X >= len(f.field[0]) || pos.Y < 0 || pos.Y >= len(f.field) {
		return true
	}
	return f.field[pos.Y][pos.X] == '.'
}

func (f *Field) set(pos Vector, b byte) {
	f.field[pos.Y][pos.X] = b
}

func (f Field) Copy() Field {
	field := make([][]byte, len(f.field))
	for i := range f.field {
		field[i] = make([]byte, len(f.field[i]))
		copy(field[i], f.field[i])
	}
	return Field{
		xMin:  f.xMin,
		xMax:  f.xMax,
		yMin:  f.yMin,
		yMax:  f.yMax,
		field: field,
	}
}

// SimulateSandDrop simulates a sand drop from the top of the field and returns true if the sand drop stays on the field
func (f *Field) SimulateSandDrop() bool {
	sand := Vector{500 - f.xMin, 0}

	for sand.X >= 0 && sand.X < len(f.field[0]) && sand.Y < len(f.field) {
		if f.isAir(Vector{sand.X, sand.Y + 1}) {
			sand.Y++
		} else if f.isAir(Vector{sand.X - 1, sand.Y + 1}) {
			sand.Y++
			sand.X--
		} else if f.isAir(Vector{sand.X + 1, sand.Y + 1}) {
			sand.Y++
			sand.X++
		} else {
			// edge case where the sand is on the start of the field
			if !f.isAir(sand) {
				return false
			}
			f.set(sand, 'o')
			return true
		}
	}

	return false
}

type Vector struct {
	X, Y int
}

func (v Vector) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

type Line struct {
	From, To Vector
}

func NewLine(from, to Vector) Line {
	// lines are  always horizontal or vertical
	if from.X > to.X || from.Y > to.Y {
		from, to = to, from
	}
	return Line{from, to}
}

func (l Line) String() string {
	return fmt.Sprintf("%v -> %v", l.From, l.To)
}

func parse(reader io.Reader) Field {
	f := NewField()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		txt := scanner.Text()
		lines := parseLine(txt)
		f.Add(lines)
	}

	return f.Compile()
}

func parseLine(s string) []Line {
	parts := strings.Split(s, " -> ")
	vecs := list.Transform(parts, func(t string) Vector {
		var x, y int
		_, _ = fmt.Sscanf(t, "%d,%d", &x, &y)
		return Vector{x, y}
	})

	lines := make([]Line, 0, len(vecs)-1)
	for i := 0; i < len(vecs)-1; i++ {
		lines = append(lines, NewLine(vecs[i], vecs[i+1]))
	}

	return lines
}

// --- helpers

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

// -- execution

func main() {
	reader, err := aoc.Get(14)
	if err != nil {
		panic(err)
	}
	f := parse(reader)

	r1 := solve1(f.Copy())
	fmt.Println("Part 1:", r1)

	r2 := solve2(f.Copy())
	fmt.Println("Part 2:", r2)
}

func solve1(f Field) int {
	var total int
	for f.SimulateSandDrop() {
		total++
	}
	fmt.Println(f)
	return total
}

func solve2(f Field) int {

	// expand the field to the left and right
	newWidth := f.yMax + 2

	xNorm := 500 - f.xMin
	xMinAdd := newWidth - xNorm
	xMin := f.xMin - xMinAdd

	xMax := xMin + (2 * newWidth)

	padLeft := bytes.Repeat([]byte{'.'}, f.xMin-xMin)
	padRight := bytes.Repeat([]byte{'.'}, xMax-f.xMax)

	for i := range f.field {
		f.field[i] = append(padLeft, f.field[i]...)
		f.field[i] = append(f.field[i], padRight...)
	}

	airgap := bytes.Repeat([]byte{'.'}, 2*newWidth+1)
	f.field = append(f.field, airgap)

	floor := bytes.Repeat([]byte{'#'}, 2*newWidth+1)
	f.field = append(f.field, floor)

	f.xMin = xMin
	f.xMax = xMax
	f.yMax = len(f.field) - 1

	var total int
	for f.SimulateSandDrop() {
		total++
	}
	return total
}

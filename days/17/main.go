package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strings"
)

func main() {
	r, err := aoc.Get(17)
	if err != nil {
		panic(err)
	}

	input, err := parse(r)
	if err != nil {
		panic(err)
	}

	heights := solve2(input, 2022)
	fmt.Println("Part 1:", heights[0])

	const rocks = 1000000000000

	initialRocks, initialHeight := 1729, 2666    // measured by hand from the output of part 1
	rockerPerCycle, heightPerCycle := 1740, 2681 // measured by hand from the output of part 1

	resultHeight := initialHeight

	fullCycles := (rocks - initialRocks) / rockerPerCycle

	resultHeight += fullCycles * heightPerCycle

	remainingRocks := (rocks - initialRocks) % rockerPerCycle

	vals := solve2(input, initialRocks+remainingRocks)
	fmt.Println(vals)

	diffH := vals[0] - initialHeight
	fmt.Println("diffH", diffH)

	resultHeight += diffH

	fmt.Println("Part 2:", resultHeight)

}

type Jets struct {
	data   []byte
	offset int
}

func (j *Jets) Next() {
	j.offset++
}

func (j *Jets) Move() Vec2 {
	if j.data[j.offset%len(j.data)] == '<' {
		return Vec2{X: -1}
	}
	return Vec2{X: 1}
}

func parse(reader io.Reader) (Jets, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return Jets{}, err
	}

	return Jets{data: bytes.Split(b, []byte{'\n'})[0]}, nil
}

var (
	HLineString = `####`
	PlusString  = `.#.
###
.#.`
	CornersString = `..#
..#
###`

	VLineString = `#
#
#
#`

	SquareString = `##
##`
)

var (
	VerticalLine   = NewShape(VLineString)
	HorizontalLine = NewShape(HLineString)
	Plus           = NewShape(PlusString)
	Corner         = NewShape(CornersString)
	Square         = NewShape(SquareString)
)

var shapes = []Shape{HorizontalLine, Plus, Corner, VerticalLine, Square}

type Field struct {
	height int
	width  int
	blocks []Block
}

func (f Field) CanMove(shape Shape, pos Vec2) bool {
	block := Block{
		Position: pos,
		Shape:    shape,
	}
	if block.MinX() < 0 {
		return false
	}
	if block.MaxX() > f.width-1 {
		return false
	}

	if block.MinY() < 0 {
		return false
	}

	return !f.HasCollisonWithBlocks(block)
}

func (f Field) HasCollisonWithBlocks(block Block) bool {
	for i := len(f.blocks) - 1; i >= 0; i-- {
		if f.blocks[i].Collides(block) {
			return true
		}
	}
	return false
}

func (f *Field) Add(shape Shape, pos Vec2) {
	block := Block{
		Position: pos,
		Shape:    shape,
	}
	f.height = max(f.height, 1+block.MaxY())
	f.blocks = append(f.blocks, block)
}

func (f Field) Copy() Field {
	cp := Field{
		height: f.height,
		width:  f.width,
		blocks: make([]Block, len(f.blocks)),
	}

	for i, b := range f.blocks {
		cp.blocks[i] = b
	}

	return cp
}

func (f Field) String() string {

	allPts := make([]Vec2, 0, len(f.blocks)*5)

	for _, b := range f.blocks {
		allPts = append(allPts, b.Point()...)
	}

	var maxY int
	for _, p := range allPts {
		maxY = max(maxY, p.Y)
	}

	bts := make([][]byte, maxY+1)
	for i := range bts {
		bts[i] = bytes.Repeat([]byte{'.'}, f.width)
	}

	for _, block := range f.blocks {
		pts := block.Point()
		for _, p := range pts {
			y := maxY - p.Y
			bts[y][p.X] = '@'
		}
	}

	var sb strings.Builder
	for _, l := range bts {
		sb.WriteByte('|')
		sb.Write(l)
		sb.WriteByte('|')
		sb.WriteByte('\n')
	}

	sb.WriteRune('+')
	sb.Write(bytes.Repeat([]byte{'-'}, f.width))
	sb.WriteRune('+')

	return sb.String()
}

type Bounds struct {
	minX, maxX, minY, maxY int
}

type Block struct {
	Position Vec2
	Shape    Shape
}

func (b Block) MinY() int {
	return b.Position.Y
}

func (b Block) MaxY() int {
	return b.Position.Y + b.Shape.height - 1
}

func (b Block) MinX() int {
	return b.Position.X
}

func (b Block) MaxX() int {
	return b.Position.X + b.Shape.width - 1
}

func (b Block) Bounds() Bounds {
	return Bounds{
		minX: b.MinX(),
		maxX: b.MaxX(),
		minY: b.MinY(),
		maxY: b.MaxY(),
	}
}

func (b Block) Collides(other Block) bool {
	// check X overlap
	if other.MaxX() < b.MinX() || other.MinX() > b.MaxX() ||
		other.MaxY() < b.MinY() || other.MinY() > b.MaxY() {
		return false
	}

	ps1, ps2 := b.Point(), other.Point()
	for _, p1 := range ps1 {
		for _, p2 := range ps2 {
			if p1.X == p2.X && p1.Y == p2.Y {
				return true
			}
		}
	}

	return false
}

func (b Block) Point() []Vec2 {
	points := make([]Vec2, 0, len(b.Shape.fields))
	for _, p := range b.Shape.fields {
		points = append(points, b.Position.Add(p.X, p.Y))
	}
	return points
}

type Vec2 struct {
	X, Y int
}

func (v Vec2) Add(x, y int) Vec2 {
	return Vec2{
		X: v.X + x,
		Y: v.Y + y,
	}
}

type Shape struct {
	height int
	width  int
	fields []Vec2
}

func NewShape(s string) Shape {
	numFields := strings.Count(s, "#")
	lines := strings.Split(s, "\n")
	width := len(lines[0])
	height := len(lines)

	fields := make([]Vec2, 0, numFields)

	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				fields = append(fields, Vec2{X: x, Y: height - 1 - y})
			}
		}
	}

	return Shape{
		width:  width,
		height: height,
		fields: fields,
	}
}

func (s Shape) Width() int {
	return s.width
}

func (s Shape) Height() int {
	return s.height
}

type Mover struct {
	isHorizontal bool
	stream       []byte
	offset       int
}

func NewMover(stream Jets) *Mover {
	return &Mover{
		isHorizontal: true,
		stream:       stream.data,
	}
}

func (m *Mover) Reset() {
	m.isHorizontal = false
}

func (m *Mover) Next() (Vec2, bool) {
	m.isHorizontal = !m.isHorizontal

	if m.isHorizontal {
		c := m.stream[m.offset%len(m.stream)]
		m.offset++
		if c == '<' {
			return Vec2{X: -1}, true
		}

		return Vec2{X: 1}, true
	}

	return Vec2{Y: -1}, false
}

func solve1(input Jets, blocks int) int {
	field := Field{
		width:  7,
		blocks: make([]Block, 0, blocks),
	}

	defaultOffset := Vec2{
		X: 2,
		Y: 3,
	}

	mover := NewMover(input)

	heights := make([]int, blocks)

	for i := 0; i < blocks; i++ {
		shape := shapes[i%len(shapes)]
		position := defaultOffset.Add(0, field.height)
		mover.Reset()

		for {
			move, isHorizontal := mover.Next()

			next := position.Add(move.X, move.Y)

			if field.CanMove(shape, next) {
				position = next
			} else if !isHorizontal {
				break
			}
		}

		field.Add(shape, position)
		heights[i] = field.height
	}

	return field.height
}

func solve2(input Jets, blocks int) []int {
	field := Field{
		width:  7,
		blocks: make([]Block, 0, blocks),
	}

	defaultOffset := Vec2{
		X: 2,
		Y: 3,
	}

	mover := NewMover(input)

	heights := make([]int, blocks)

	var cycle int
	var lastFallen int
	var lastHeight int

	iFallen, iHeight := 0, 0
	rps, hps := 0, 0

	for i := 0; i < blocks; i++ {
		shape := shapes[i%len(shapes)]
		position := defaultOffset.Add(0, field.height)
		mover.Reset()

		for {
			move, isHorizontal := mover.Next()

			if mover.offset%len(mover.stream) == 0 && cycle%2 == 0 {
				rps = len(field.blocks) - lastFallen
				hps = field.height - lastHeight
				if cycle == 0 {
					iFallen, iHeight = rps, hps
				}
				fmt.Println("CYCLE", cycle)
				fmt.Println(rps, hps)

				lastHeight = field.height
				lastFallen = len(field.blocks)
				cycle++
			}

			next := position.Add(move.X, move.Y)

			if field.CanMove(shape, next) {
				position = next
			} else if !isHorizontal {
				break
			}
		}

		field.Add(shape, position)
		heights[i] = field.height
	}

	return []int{field.height, iFallen, iHeight, rps, hps}
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

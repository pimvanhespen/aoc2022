package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

func main() {
	reader, err := aoc.Get(9)
	if err != nil {
		panic(err)
	}

	input, err := parse(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(input))
	fmt.Println("Part 2:", solve2(input))

}

var (
	Unknown = Vector{0, 0}
	Up      = Vector{0, -1}
	Right   = Vector{1, 0}
	Down    = Vector{0, 1}
	Left    = Vector{-1, 0}
)

type Move struct {
	Direction Vector
	Steps     int
}

type Input []Move

func parse(reader io.Reader) (Input, error) {

	var input Input

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		var c byte
		var d int
		_, _ = fmt.Sscanf(text, "%c %d", &c, &d)

		var dir Vector
		switch c {
		case 'U':
			dir = Up
		case 'R':
			dir = Right
		case 'D':
			dir = Down
		case 'L':
			dir = Left
		}

		input = append(input, Move{
			Direction: dir,
			Steps:     d,
		})
	}

	return input, nil
}

// Vector represents a 2D vector.
type Vector struct {
	X int
	Y int
}

// Add adds the other vector to the current vector.
// e.g.: Vector{1, 2}.Add(Vector{3, 4}) == Vector{4, 6}
func (v Vector) Add(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub subtracts the other vector from the current vector.
// e.g.: Vector{1, 2}.Sub(Vector{3, 4}) == Vector{-2, -2}
func (v Vector) Sub(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Distance returns the Manhattan distance between two vectors.
// e.g.: Vector{1, 2}.Distance(Vector{3, 4}) == 4
func (v Vector) Distance(other Vector) int {
	return abs(v.X-other.X) + abs(v.Y-other.Y)
}

// Normalize returns a vector with both x and y set to either -1, 0 or 1.
// e.g.: Vector{-10, 10}.Normalize() == Vector{-1, 1}
func (v Vector) Normalize() Vector {
	return Vector{limit(v.X), limit(v.Y)}
}

func Follow(position, target Vector) Vector {
	diff := target.Sub(position)
	dist := position.Distance(target)
	norm := diff.Normalize()

	if diff.X != 0 && diff.Y != 0 { // keep this if separate from subsequent checks!
		// diagonal
		if dist > 2 {
			return position.Add(norm)
		}

	} else if diff.X != 0 && dist > 1 {
		// move horizontally
		return position.Add(norm)

	} else if diff.Y != 0 && dist > 1 {
		// move vertically
		return position.Add(norm)
	}

	return position // no move
}

func solve1(input Input) int {
	return solve(input, 2)
}

func solve2(input Input) int {
	return solve(input, 10)
}

func solve(input Input, knots int) int {

	//tailMoves := make(map[Vector]struct{}) // tail moves <- used in original solution

	visited := make([]map[Vector]struct{}, knots)
	for i := range visited {
		visited[i] = make(map[Vector]struct{})
		visited[i][Vector{0, 0}] = struct{}{} // start at 0,0
	}

	registerMove := func(position Vector, tail int) {
		if _, ok := visited[tail][position]; !ok {
			visited[tail][position] = struct{}{}
		}
	}

	rope := make([]Vector, knots)

	for _, move := range input {
		for i := 0; i < move.Steps; i++ {

			// head moves
			rope[0] = rope[0].Add(move.Direction)
			registerMove(rope[0], 0)

			// knots follow their predecessor
			for knot := 1; knot < len(rope); knot++ {
				rope[knot] = Follow(rope[knot], rope[knot-1])
				registerMove(rope[knot], knot)
			}
		}
	}

	// store outputs
	floor := renderFloor(visited[len(visited)-1], rope)
	aoc.SaveOutput(9, fmt.Sprintf("floor_%d.txt", knots), func(writer io.Writer) error {
		_, err := writer.Write(floor)
		return err
	})

	img := createImage(visited)
	filename := fmt.Sprintf("image_%d.png", knots)
	aoc.SaveOutput(9, filename, func(writer io.Writer) error {
		return png.Encode(writer, img)
	})

	return len(visited[len(visited)-1])
}

//														//
// -- rendering -- NOT A PART OF THE ACTUAL SOLUTION -- //
//														//

// renderFloor prints the floor with all the tiles visited marked with '#'.
// The head and tail of the rope are marked with 'H' and 'T' respectively.
// The start position is marked with 'S'.
// The intermediate knots of the rope are marked with '0-9a-z'
func renderFloor(tailHistory map[Vector]struct{}, rope []Vector) []byte {
	var yMin, yMax, xMin, xMax int

	updateLimits := func(v Vector) {
		xMin = min(xMin, v.X)
		xMax = max(xMax, v.X)
		yMin = min(yMin, v.Y)
		yMax = max(yMax, v.Y)
	}

	for position := range tailHistory {
		updateLimits(position)
	}

	for _, position := range rope {
		updateLimits(position)
	}

	width := abs(xMax) + abs(xMin) + 1  // +1 for the start position (0,0)
	height := abs(yMax) + abs(yMin) + 1 // +1 for the start position (0,0)

	floor := make([][]byte, height)
	for i := range floor {
		floor[i] = bytes.Repeat([]byte{'.'}, width)
	}

	setTile := func(v Vector, c byte) {
		floor[v.Y-yMin][v.X-xMin] = c
	}

	for position := range tailHistory {
		setTile(position, '#')
	}

	setTile(Vector{0, 0}, 's')

	for i := 1; i < len(rope)-1; i++ {
		setTile(rope[i], byte('0'+i))
	}

	// head and tail
	setTile(rope[0], 'H')
	setTile(rope[len(rope)-1], 'T')

	// reverse output order so that the output aligns with the coords.
	// By default, (0,0) is in the top left corner, because going up in an array is y-- rather than y++.
	for i := 0; i < len(floor)/2; i++ {
		opposite := len(floor) - i - 1
		floor[i], floor[opposite] = floor[opposite], floor[i]
	}

	return bytes.Join(floor, []byte{'\n'})
}

func createImage(history []map[Vector]struct{}) image.Image {
	var yMin, yMax, xMin, xMax int
	for _, tail := range history {
		for position := range tail {
			xMin = min(xMin, position.X)
			xMax = max(xMax, position.X)
			yMin = min(yMin, position.Y)
			yMax = max(yMax, position.Y)
		}
	}

	width := abs(xMax) + abs(xMin) + 1 + 2  // +1 for the start position (0,0)
	height := abs(yMax) + abs(yMin) + 1 + 2 // +1 for the start position (0,0)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width*height; i++ {
		img.Set(i%width, i/width, color.RGBA{0, 0, 0, 255})
	}

	// create a palette
	palette := color.Palette{
		color.RGBA{0xff, 0x00, 0x00, 0xFF}, // red
		color.RGBA{0xff, 0x45, 0x00, 0xFF}, // orange-red
		color.RGBA{0xff, 0xA5, 0x00, 0xFF}, // orange
		color.RGBA{0xff, 0xff, 0x00, 0xFF}, // yellow
		color.RGBA{0xAD, 0xFF, 0x2F, 0xFF}, // green-yellow
		color.RGBA{0x00, 0xff, 0x00, 0xFF}, // green
		color.RGBA{0x3C, 0xB3, 0x71, 0xFF}, // medium sea green
		color.RGBA{0x00, 0xff, 0xff, 0xFF}, // cyan
		color.RGBA{0x00, 0x00, 0xff, 0xFF}, // blue
		color.RGBA{0x4B, 0x00, 0x82, 0xFF}, // indigo
		color.RGBA{0x8A, 0x2B, 0xE2, 0xFF}, // blue-violet
		color.RGBA{0xff, 0x00, 0xff, 0xFF}, // magenta
		color.RGBA{0xff, 0x00, 0x7F, 0xFF}, // deep pink
	}

	setTile := func(v Vector, c color.Color, offset int) {
		prev := img.At(v.X-xMin+1, v.Y-yMin+1)
		blended := blendColors(prev, c, 1.0/float64(offset+2))
		img.Set(v.X-xMin+1, v.Y-yMin+1, blended)
	}

	for i, tail := range history {
		if i >= len(palette) {
			panic("not enough colors in palette")
		}

		c := palette[i]
		if i == len(history)-1 {
			c = color.RGBA{0xff, 0xff, 0xff, 0xff}
		}

		for position := range tail {
			setTile(position, c, 0)
			for _, offset := range []Vector{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				setTile(position.Add(offset), c, 1)
			}
		}
	}

	return img
}

func reduceRGB(c color.Color, f float64) color.Color {
	_, _, _, a := c.RGBA()
	return blendColors(c, color.RGBA{0, 0, 0, uint8(a / 0xFF)}, f)
}

func reduceAlpha(c color.Color, f float64) color.Color {
	r, g, b, a := c.RGBA()

	return color.RGBA{
		R: uint8(r / 0xff),
		G: uint8(g / 0xff),
		B: uint8(b / 0xff),
		A: uint8(float64(a/0xff) * f),
	}
}

func limit(n int) int {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func blendColorValue(a, b, t float64) float64 {
	// https://stackoverflow.com/questions/726549/algorithm-for-additive-color-mixing-for-rgb-values
	// sqrt((1-t)*a ^ 2 + t*b ^ 2)
	return math.Sqrt((1-t)*a*a + t*b*b)
}

func blendAlphaValue(a, b, t float64) float64 {
	return (1-t)*a + t*b
}

func normalize(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := c.RGBA()
	return float64(r) / 0xffff,
		float64(b) / 0xffff,
		float64(g) / 0xffff,
		float64(a) / 0xffff
}

func denormalize(r, g, b, a float64) color.RGBA {
	return color.RGBA{
		R: uint8(r * 0xff),
		G: uint8(g * 0xff),
		B: uint8(b * 0xff),
		A: uint8(a * 0xff),
	}
}

func blendColors(c1, c2 color.Color, t float64) color.RGBA {
	r1, g1, b1, a1 := normalize(c1)
	r2, g2, b2, a2 := normalize(c2)

	return denormalize(
		blendColorValue(r1, r2, t),
		blendColorValue(g1, g2, t),
		blendColorValue(b1, b2, t),
		blendAlphaValue(a1, a2, t),
	)
}

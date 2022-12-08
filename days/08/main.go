package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
)

func main() {
	rc, err := aoc.Get(8)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	field, err := parse(rc)
	if err != nil {
		panic(err)
	}

	// Part 1
	r1, err := solve1(field)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", r1)

	r2, err := solve2(field)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", r2)
}

func parse(rc io.Reader) ([][]byte, error) {

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var field [][]byte

	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		field = append(field, line)
	}

	return field, nil
}

func solve1(field [][]byte) (int, error) {

	height := len(field)
	width := len(field[0])

	var visible int

	for i := 0; i < width*height; i++ {
		x := i % width
		y := i / width

		if isVisible(field, x, y) {
			visible++
		}
	}

	return visible, nil
}

func solve2(field [][]byte) (int, error) {

	height := len(field)
	width := len(field[0])

	var highscore int

	var score int

	for i := 0; i < width*height; i++ {
		x := i % width
		y := i / width

		score = scenicScore(field, x, y)

		if score > highscore {
			highscore = score
		}
	}

	return highscore, nil
}

func isVisible(field [][]byte, xPos, yPos int) bool {

	if xPos <= 0 || xPos >= len(field[0])-1 {
		return true
	}

	if yPos <= 0 || yPos >= len(field)-1 {
		return true
	}

	return !isBlocked(field, xPos, yPos, 0, -1) ||
		!isBlocked(field, xPos, yPos, 0, 1) ||
		!isBlocked(field, xPos, yPos, -1, 0) ||
		!isBlocked(field, xPos, yPos, 1, 0)
}

func isBlocked(field [][]byte, x, y, dX, dY int) bool {
	if dX == 0 && dY == 0 {
		return false
	}

	sx, sy := x+dX, y+dY

	var min, height uint8

	min = field[y][x] - '0'

	for sx >= 0 && sx < len(field[0]) && sy >= 0 && sy < len(field) {
		height = field[sy][sx] - '0'

		if height >= min {
			return true
		}

		sx += dX
		sy += dY
	}
	return false
}

func scenicScore(field [][]byte, x, y int) int {
	return scenicScoreInDirection(field, x, y, 0, -1) *
		scenicScoreInDirection(field, x, y, 0, 1) *
		scenicScoreInDirection(field, x, y, -1, 0) *
		scenicScoreInDirection(field, x, y, 1, 0)
}

func scenicScoreInDirection(field [][]byte, x, y, dX, dY int) int {
	if dX == 0 && dY == 0 {
		return 0
	}

	sx, sy := x+dX, y+dY

	var min, height uint8

	min = field[y][x] - '0'

	var score int

	for sx >= 0 && sx < len(field[0]) && sy >= 0 && sy < len(field) {
		height = field[sy][sx] - '0'

		score++

		if height >= min {
			return score
		}

		sx += dX
		sy += dY
	}
	return score
}

// -- alternative solutions --

func solve1_visibility_matrix(field [][]byte) (int, error) {
	vis := make([][]bool, len(field))
	for i := 0; i < len(field); i++ {
		vis[i] = make([]bool, len(field[0]))
	}

	var max, cur uint8

	for y := 0; y < len(field); y++ {
		if y == 0 || y == len(field)-1 {
			for x := 0; x < len(field[0]); x++ {
				vis[y][x] = true
			}
		} else {
			vis[y][0] = true
			vis[y][len(field[0])-1] = true
		}
	}

	for y := 1; y < len(field)-1; y++ {

		max = field[y][0]

		// left to right
		for x := 1; x < len(field[0])-1; x++ {
			cur = field[y][x]

			if cur > max {
				max = cur
				vis[y][x] = true
			}
		}

		max = field[y][len(field[0])-1]

		// right to left
		for x := len(field[0]) - 2; x > 0; x-- {
			cur = field[y][x]

			if cur > max {
				max = cur
				vis[y][x] = true
			}
		}
	}

	for x := 1; x < len(field[0])-1; x++ {
		max = field[0][x]

		// top to bottom
		for y := 1; y < len(field)-1; y++ {
			cur = field[y][x]

			if cur > max {
				max = cur
				vis[y][x] = true
			}
		}

		max = field[len(field)-1][x]

		// bottom to top
		for y := len(field) - 2; y > 0; y-- {
			cur = field[y][x]

			if cur > max {
				max = cur
				vis[y][x] = true
			}
		}
	}

	//printVisibility(vis)

	visible := 2 * (len(field) + len(field[0]) - 2)

	for y := 1; y < len(field)-1; y++ {
		for x := 1; x < len(field[0])-1; x++ {
			if vis[y][x] {
				visible++
			}
		}
	}

	return visible, nil
}

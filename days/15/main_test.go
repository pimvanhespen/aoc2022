package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"
)

var testInput = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`

func TestSolve1(t *testing.T) {
	const want = 26
	sensors := parse(strings.NewReader(testInput))

	printMap(sensors)

	got := solve1(sensors, 10)
	if got != want {
		t.Errorf("solve1() = %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	const want = 56000011
	sensors := parse(strings.NewReader(testInput))

	got := solve2(sensors, 0, 20)
	if got != want {
		t.Errorf("solve2() = %d, want %d", got, want)
	}
}

func printMap(sensors []Sensor) {
	minX := math.MaxInt
	maxX := 0
	minY := math.MaxInt
	maxY := 0

	for _, s := range sensors {
		dist := ManhattanDistance(s.Position, s.NearestBeacon)

		minX = Min(minX, s.Position.X-dist)
		maxX = Max(maxX, s.Position.X+dist)
		minY = Min(minY, s.Position.Y-dist)
		maxY = Max(maxY, s.Position.Y+dist)
	}

	m := make([][]byte, maxY-minY+1)
	for i := range m {
		m[i] = bytes.Repeat([]byte{'.'}, maxX-minX+1)
	}

	for yMap := 0; yMap < len(m); yMap++ {
		yOffset := yMap + minY
		for _, s := range sensors {
			minmax, ok := s.CalcXAxisCoverage(yOffset)
			if !ok {
				continue
			}
			for x := minmax[Low]; x <= minmax[High]; x++ {
				xMap := x - minX
				m[yMap][xMap] = '#'
			}
		}
	}

	for _, s := range sensors {
		m[s.Position.Y-minY][s.Position.X-minX] = 'S'
		m[s.NearestBeacon.Y-minY][s.NearestBeacon.X-minX] = 'B'
	}

	header := make([][]byte, 3)
	for i := range header {
		header[i] = bytes.Repeat([]byte{' '}, maxX-minX+1+4)
	}

	for i := 4; i < len(header[0]); i++ {
		if !(i == 4 || i == len(header[0])-1 || (i+minX)%5 == 0) {
			continue
		}

		s := fmt.Sprintf("%3d", i+minX)
		header[0][i] = s[0]
		header[1][i] = s[1]
		header[2][i] = s[2]
	}

	m = append(header, m...)

	for n := range m[3:] {
		bts := []byte("    ")

		y := n + minY
		if n == 0 || y == len(m)-1 || (y)%5 == 0 || y == 0 {
			s := fmt.Sprintf("%3d", y)
			bts[0] = s[0]
			bts[1] = s[1]
			bts[2] = s[2]
		}

		m[n+3] = append(bts, m[n+3]...)
	}

	fmt.Println(string(bytes.Join(m, []byte{'\n'})))
}

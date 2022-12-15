package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"sort"
)

const (
	Low  = 0
	High = 1
)

type Vector struct {
	X, Y int
}

func (v Vector) Add(x, y int) Vector {
	return Vector{
		X: v.X + x,
		Y: v.Y + y,
	}
}

type Sensor struct {
	Position      Vector
	NearestBeacon Vector
}

func (s Sensor) CalcXAxisCoverage(y int) ([2]int, bool) {
	dist := ManhattanDistance(s.Position, s.NearestBeacon)
	distToY := Abs(s.Position.Y - y)

	radius := dist - distToY

	if radius < 0 {
		return [2]int{}, false
	}

	return [2]int{s.Position.X - radius, s.Position.X + radius}, true
}

func ManhattanDistance(a, b Vector) int {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	r, err := aoc.Get(15)
	if err != nil {
		panic(err)
	}

	sensors := parse(r)

	part1 := solve1(sensors, 2_000_000)
	fmt.Println("Part1:", part1)

	part2 := solve2(sensors, 0, 4_000_000)
	fmt.Println("Part2:", part2)
}

func parse(reader io.Reader) []Sensor {

	var sensors []Sensor

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		var b Sensor
		_, _ = fmt.Sscanf(text, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &b.Position.X, &b.Position.Y, &b.NearestBeacon.X, &b.NearestBeacon.Y)
		sensors = append(sensors, b)
	}

	return sensors
}

func solve1(sensors []Sensor, y int) int {

	var coverages [][2]int
	beacons := map[Vector]struct{}{}

	for _, s := range sensors {

		rng, ok := s.CalcXAxisCoverage(y)
		if !ok {
			continue
		}

		coverages = append(coverages, rng)

		if s.NearestBeacon.Y == y {
			beacons[s.NearestBeacon] = struct{}{}
		}
	}

	coverages = merge(coverages)

	tilesCovered := 0

	for _, c := range coverages {
		tilesCovered += 1 + c[High] - c[Low] // 1 + width
	}

	return tilesCovered - len(beacons) // tiles with beacons don't count
}

func merge(coverages [][2]int) [][2]int {
	sort.Slice(coverages, func(i, j int) bool {
		return coverages[i][Low] < coverages[j][Low]
	})

	for {
		before := len(coverages)

		for i := len(coverages) - 1; i > 0; i-- {
			if coverages[i][Low] > coverages[i-1][High]+1 {
				continue
			}

			coverages[i-1][Low] = Min(coverages[i-1][Low], coverages[i][Low])
			coverages[i-1][High] = Max(coverages[i-1][High], coverages[i][High])

			coverages = append(coverages[:i], coverages[i+1:]...)
		}

		if len(coverages) == before {
			break
		}
	}

	return coverages
}

func solve2(sensors []Sensor, minC, maxC int) int {

	var coverages [][2]int

	var x, y int

	for y = minC; y <= maxC; y++ {
		coverages = coverages[:0]

		for _, s := range sensors {

			minmax, ok := s.CalcXAxisCoverage(y)
			if !ok || minmax[High] < minC || minmax[Low] > maxC {
				continue
			}

			coverages = append(coverages, minmax)
		}

		coverages = merge(coverages)

		if len(coverages) <= 0 {
			panic("no coverage")
		}

		if len(coverages) >= 3 {
			panic("too many coverages")
		}

		if len(coverages) == 1 {
			if coverages[Low][Low] <= minC && coverages[Low][High] >= maxC {
				continue
			}

			if coverages[Low][Low] > minC && coverages[Low][High] < maxC {
				panic("coverage is too small")
			}

			// we have a single coverage that does not cover the whole range
			if coverages[Low][Low] == minC {
				x = maxC
			} else {
				x = minC
			}
			break
		}

		if len(coverages) == 2 {
			x = coverages[Low][High] + 1
			break
		}
	}

	return y + x*4_000_000
}

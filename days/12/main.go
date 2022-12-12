package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/queue"
	"io"
	"math"
)

func main() {
	r, err := aoc.Get(12)
	if err != nil {
		panic(err)
	}

	bts, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes.Count(bts, []byte("a")))

	c, b, e := parse(bts)

	// Part 1
	p1 := solve1(c, b, e)
	fmt.Println("Part 1:", p1)

	p2 := solve2(c, e)
	fmt.Println("Part 2:", p2)
}

func parse(data []byte) (Chart, Tile, Tile) {
	data = bytes.TrimSpace(data)
	rows := bytes.Split(data, []byte{'\n'})

	chart := make([][]Tile, len(rows))

	var start, end Tile

	for y, row := range rows {
		chart[y] = make([]Tile, len(row))
		for x, b := range row {

			t := Tile{
				X:         x,
				Y:         y,
				Elevation: b,
			}

			if t.Elevation == 'S' {
				t.Elevation = 'a'
				start = t
			} else if t.Elevation == 'E' {
				t.Elevation = 'z'
				end = t
			}

			chart[y][x] = t
		}
	}

	return Chart(chart), start, end
}

// -- Solution --

func solve1(chart Chart, start, end Tile) int {

	heuristic := func(t Tile) int {
		return ManhattanDistance(t, end)
	}

	cost := func(from, to Tile) (int, bool) {
		if (from.Elevation + 1) >= to.Elevation {
			return 1, true
		}

		return 0, false
	}

	neighbours := func(t Tile) []Tile {
		return chart.GetDirectNeighbours(t.X, t.Y)
	}

	path := AStar(start, end, heuristic, cost, neighbours)
	if path == nil {
		panic("no path found")
	}
	printRoute(path)
	return len(path[1:])
}

func solve2(chart Chart, end Tile) int {
	heuristic := func(t Tile) int {
		return ManhattanDistance(t, end)
	}

	cost := func(from, to Tile) (int, bool) {
		if (from.Elevation + 1) >= to.Elevation {
			return 1, true
		}

		return 0, false
	}

	neighbours := func(t Tile) []Tile {
		return chart.GetDirectNeighbours(t.X, t.Y)
	}

	// Find clusters of 'a' so that we can skip unreachable clusters later on
	clusters := chart.Clusters('a')
	if len(clusters) == 0 {
		panic("no options")
	}

	minimum := math.MaxInt
	var shortest []Tile // for printing the output

outer:
	for _, cluster := range clusters {
		for _, tile := range cluster {
			route := AStar(tile, end, heuristic, cost, neighbours)
			if len(route) == 0 {
				// if one tile in the cluster cannot reach the top, none can!
				continue outer
			}

			steps := len(route) - 1

			if steps < minimum {
				shortest = route
				minimum = steps
			}
		}
	}

	printRoute(shortest)
	return minimum
}

type Tile struct {
	X, Y      int
	Elevation byte
}

func (t Tile) Equals(other Tile) bool {
	return t.X == other.X && t.Y == other.Y
}

type Chart [][]Tile

func (m Chart) Get(x, y int) (Tile, bool) {
	if x < 0 || y < 0 || x >= len(m[0]) || y >= len(m) {
		return Tile{}, false
	}

	return m[y][x], true
}

func (m Chart) GetDirectNeighbours(x, y int) []Tile {
	tiles := make([]Tile, 0, 4)
	if t, ok := m.Get(x-1, y); ok {
		tiles = append(tiles, t)
	}
	if t, ok := m.Get(x+1, y); ok {
		tiles = append(tiles, t)
	}
	if t, ok := m.Get(x, y-1); ok {
		tiles = append(tiles, t)
	}
	if t, ok := m.Get(x, y+1); ok {
		tiles = append(tiles, t)
	}
	return tiles
}

func (m Chart) Clusters(c byte) [][]Tile {

	var clusters [][]Tile

	cache := map[Tile]struct{}{}

	for _, row := range m {
		for _, tile := range row {
			if tile.Elevation != c {
				continue
			}

			if _, ok := cache[tile]; ok {
				continue
			}

			clusters = append(clusters, getClusterRecursive(m, cache, tile))
		}
	}

	return clusters
}

func getClusterRecursive(m Chart, cache map[Tile]struct{}, t Tile) []Tile {
	if _, ok := cache[t]; ok {
		return nil
	}

	cluster := []Tile{t}
	cache[t] = struct{}{}

	for _, n := range m.GetDirectNeighbours(t.X, t.Y) {
		if n.Elevation == t.Elevation {
			cluster = append(cluster, getClusterRecursive(m, cache, n)...)
		}
	}

	return cluster
}

func ManhattanDistance(a, b Tile) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// A* Algorithm

type HeuristicFunc func(Tile) int
type CostFunc func(Tile, Tile) (int, bool)
type NeighboursFunc func(Tile) []Tile

func AStar(start, end Tile, h HeuristicFunc, costFn CostFunc, neighbours NeighboursFunc) []Tile {

	pq := queue.NewPriority[Tile]()
	pq.Insert(start, 0)

	cameFrom := make(map[Tile]Tile)

	gScore := make(map[Tile]int)
	gScore[start] = 0

	for pq.Len() > 0 {

		current := pq.Pop()
		if current == end {
			return reconstructPath(cameFrom, current)
		}

		for _, neighbour := range neighbours(current) {

			price, allowed := costFn(current, neighbour)
			if !allowed {
				continue
			}

			tentativeGScore := gScore[current] + price

			// register if score is better, or tile has no score yet
			if gScoreNeighbour, ok := gScore[neighbour]; !ok || tentativeGScore < gScoreNeighbour {
				cameFrom[neighbour] = current
				gScore[neighbour] = tentativeGScore
				fScore := tentativeGScore + h(neighbour)
				pq.Upsert(neighbour, fScore)
			}
		}

	}

	return nil
}

func reconstructPath(cameFrom map[Tile]Tile, current Tile) []Tile {
	totalPath := []Tile{current}
	for {
		prev, ok := cameFrom[current]
		if !ok {
			break
		}

		current = prev
		totalPath = append([]Tile{prev}, totalPath...)
	}
	return totalPath
}

// ---- End of solutions //

func printRoute(path []Tile) {

	var width, height int
	for _, step := range path {
		width, height = max(width, step.X), max(height, step.Y)
	}

	chart := make([][]byte, height+1) // +1 offset as coords are zero indexed, height start at 1
	for i := range chart {
		chart[i] = bytes.Repeat([]byte{'.'}, width+1) // +1 offset as coords are zero indexed
	}

	for _, step := range path {
		chart[step.Y][step.X] = step.Elevation
	}

	fmt.Println(string(bytes.Join(chart, []byte{'\n'})))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

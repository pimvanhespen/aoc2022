package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/set"
	"io"
	"log"
	"math"
)

type Vector3D struct {
	X, Y, Z int
}

func (v Vector3D) String() string {
	return fmt.Sprintf("(%d,%d,%d)", v.X, v.Y, v.Z)
}

func (v Vector3D) Merge(other Vector3D, fn func(a, b int) int) Vector3D {
	return Vector3D{fn(v.X, other.X), fn(v.Y, other.Y), fn(v.Z, other.Z)}
}

type Block struct {
	Pos    Vector3D
	neighs []*Block
	isLava bool
}

func NewBlock(pos Vector3D) *Block {
	return &Block{Pos: pos}
}

func (b *Block) connect(n *Block) {
	b.neighs = append(b.neighs, n)
}

func (b Block) Neighbours() []*Block {
	return b.neighs
}

func (b Block) String() string {
	return fmt.Sprintf("%v; %t", b.Pos, b.isLava)
}

func main() {
	r, err := aoc.Get(18)
	if err != nil {
		log.Fatal(err)
	}
	v, err := parse(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(solve1(v))
	fmt.Println(solve2(v))
}

func parse(reader io.Reader) ([]Vector3D, error) {
	var v []Vector3D

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		var x, y, z int
		_, err := fmt.Sscanf(text, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			return nil, fmt.Errorf("parse '%s': %v", text, err)
		}
		v = append(v, Vector3D{x, y, z})
	}
	return v, nil
}

func solve1(v []Vector3D) int {
	exposed := len(v) * 6
	for i := 0; i < len(v)-1; i++ {
		for j := i + 1; j < len(v); j++ {
			if dist(v[i], v[j]) == 1 {
				exposed -= 2
			}
		}
	}
	return exposed
}

func solve2(v []Vector3D) int {

	vMin := Vector3D{math.MaxInt, math.MaxInt, math.MaxInt}
	vMax := Vector3D{math.MinInt, math.MinInt, math.MinInt}

	for i := 0; i < len(v); i++ {
		vMin = vMin.Merge(v[i], min)
		vMax = vMax.Merge(v[i], max)
	}

	// calc lower and upper bounds
	vMin = vMin.Merge(Vector3D{-1, -1, -1}, add)
	vMax = vMax.Merge(Vector3D{1, 1, 1}, add)

	// calc the size of the 3D space
	xDiff := vMin.Merge(vMax, absDiff)

	// Generate all blocks in the 3D space
	b := make([][][]*Block, xDiff.X+3)
	for i := range b {
		b[i] = make([][]*Block, xDiff.Y+3)
		for j := range b[i] {
			b[i][j] = make([]*Block, xDiff.Z+3)
			for k := range b[i][j] {
				x, y, z := vMin.X+i, vMin.Y+j, vMin.Z+k
				b[i][j][k] = NewBlock(Vector3D{x, y, z})
			}
		}
	}

	// Mark lava blocks as lava
	for i := 0; i < len(v); i++ {
		x, y, z := v[i].X-vMin.X, v[i].Y-vMin.Y, v[i].Z-vMin.Z
		b[x][y][z].isLava = true
	}

	// Register all neighbours for each block
	for x, r2 := range b {
		for y, r1 := range r2 {
			for z, block := range r1 {
				if x > 0 {
					block.connect(b[x-1][y][z])
				}
				if x < len(b)-1 {
					block.connect(b[x+1][y][z])
				}
				if y > 0 {
					block.connect(b[x][y-1][z])
				}
				if y < len(r2)-1 {
					block.connect(b[x][y+1][z])
				}
				if z > 0 {
					block.connect(b[x][y][z-1])
				}
				if z < len(r1)-1 {
					block.connect(b[x][y][z+1])
				}
			}
		}
	}

	// find all the outer lava blocks by wrapping an 'air blanket' around the structure

	lava := set.New[*Block]()
	seen := set.New[*Block]()
	todo := set.New[*Block]()

	todo.Add(b[0][0][0]) // begin with one of the outer air blocks

	// we loop through all directly connected air blocks, if a direct block is not an air block
	// it is an outer-layer lava block. We want to know which lava blocks are on the outside, and which air block are
	// connected to these lava blocks. By only checking outer air blocks, it is impossible to find any air blocks within
	// the structure.
	for !todo.IsEmpty() {
		block := todo.Pop()
		seen.Add(block)
		for _, n := range block.Neighbours() {
			if n.isLava {
				lava.Add(n)
				continue
			}
			if seen.Contains(n) {
				continue
			}
			todo.Add(n)
		}
	}

	// Not all air connected to lava blocks is part of the outside. Some air may be trapped inside.
	// So we only count air blocks connected to lava if they are in the outside-air set.
	var surface int
	for _, tile := range lava.ToSlice() {
		for _, n := range tile.Neighbours() {
			if seen.Contains(n) {
				surface++
			}
		}
	}
	return surface
}

func dist(a, b Vector3D) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func add(a, b int) int {
	return a + b
}

func absDiff(a, b int) int {
	return abs(a - b)
}

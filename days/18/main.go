package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"log"
)

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
}

type Vector3D struct {
	X, Y, Z int
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

func dist(a, b Vector3D) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

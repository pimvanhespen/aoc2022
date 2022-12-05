package main

import "github.com/pimvanhespen/aoc2022/pkg/aoc"

func main() {
	rc, err := aoc.Get(6)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	// todo
}

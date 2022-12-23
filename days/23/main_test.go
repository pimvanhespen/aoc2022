package main

import (
	"fmt"
	"strings"
	"testing"
)

var testInput = `....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..`

func TestSolve1(t *testing.T) {
	const want = 110

	field, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(field)

	got := solve1(field.Copy())

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	const want = 20

	field, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	got := solve2(field.Copy())

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

var smallTestInput = `
##
#.
..
##`

func TestSolveSmall(t *testing.T) {
	//const want = 99

	field, err := parse(strings.NewReader(smallTestInput))
	if err != nil {
		t.Fatal(err)
	}

	strat := createStrategy()

	fmt.Println(field)
	fmt.Println()
	for i := 0; i < 4; i++ {
		field.Simulate(1, strat)
		fmt.Println(field)
		fmt.Println()
	}
}

type expectation struct {
	afterRound int
	want       string
}

func TestDirections_Rotate(t *testing.T) {
	dirs := newDirections(NorthDirection, SouthDirection)
	dirs.Rotate()
	if dirs.dirs[0] != SouthDirection {
		t.Errorf("got %v, want %v", dirs.dirs[0], SouthDirection)
	}
	if dirs.dirs[1] != NorthDirection {
		t.Errorf("got %v, want %v", dirs.dirs[1], NorthDirection)
	}
}

func TestSimulateRound(t *testing.T) {

	for _, expect := range testExpectations[2:] {
		t.Run(fmt.Sprintf("after %d rounds", expect.afterRound), func(t *testing.T) {
			field, err := parse(strings.NewReader(testInput))
			if err != nil {
				t.Fatal(err)
			}

			strat := createStrategy()
			field.Simulate(expect.afterRound-1, strat)
			fmt.Println(field)
			field.Simulate(1, strat)

			got := field.String()

			wf, err := parse(strings.NewReader(expect.want))
			if err != nil {
				t.Fatal(err)
			}

			want := wf.String()

			if strings.Compare(got, want) != 0 {

				gots, wants := strings.Split(got, "\n"), strings.Split(want, "\n")

				var sb strings.Builder
				format := fmt.Sprintf("%%%ds %%%ds\n", len(wants[0]), len(gots[0]))
				_, _ = fmt.Fprintf(&sb, format, "want", "got")

				for i := 0; i < max(len(gots), len(wants)); i++ {
					var g, w string
					if i < len(gots) {
						g = gots[i]
					}
					if i < len(wants) {
						w = wants[i]
					}
					_, _ = fmt.Fprintf(&sb, format, w, g)
				}

				t.Errorf(sb.String())
			}
		})
	}
}

var testExpectations = []expectation{
	{
		afterRound: 0,
		want: `..............
..............
.......#......
.....###.#....
...#...#.#....
....#...##....
...#.###......
...##.#.##....
....#..#......
..............
..............
..............`,
	},
	{
		afterRound: 1,
		want: `..............
.......#......
.....#...#....
...#..#.#.....
.......#..#...
....#.#.##....
..#..#.#......
..#.#.#.##....
..............
....#..#......
..............
..............`,
	},
	{
		afterRound: 2,
		want: `..............
.......#......
....#.....#...
...#..#.#.....
.......#...#..
...#..#.#.....
.#...#.#.#....
..............
..#.#.#.##....
....#..#......
..............
..............`,
	},
	{
		afterRound: 3,
		want: `..............
.......#......
.....#....#...
..#..#...#....
.......#...#..
...#..#.#.....
.#..#.....#...
.......##.....
..##.#....#...
...#..........
.......#......
..............`,
	},
	{
		afterRound: 4,
		want: `..............
.......#......
......#....#..
..#...##......
...#.....#.#..
.........#....
.#...###..#...
..#......#....
....##....#...
....#.........
.......#......
..............`,
	},
	{
		afterRound: 5,
		want: `.......#......
..............
..#..#.....#..
.........#....
......##...#..
.#.#.####.....
...........#..
....##..#.....
..#...........
..........#...
....#..#......
..............`},
	{
		afterRound: 10,
		want: `.......#......
...........#..
..#.#..#......
......#.......
...#.....#..#.
.#......##....
.....##.......
..#........#..
....#.#..#....
..............
....#..#..#...
..............`,
	},
}

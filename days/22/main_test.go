package main

import (
	"reflect"
	"strings"
	"testing"
)

var testInput = `        ...#
        .#..
        #...
        ....
...#.......#
........#...
..#....#....
..........#.
        ...#....
        .....#..
        .#......
        ......#.

10R5L5R10L4R5L5`

func TestSolve1(t *testing.T) {
	field, path, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solve1(field.Copy(), path), 6032; got != want {
		t.Errorf("solve1() = %d, want %d", got, want)
	}
}

func TestParseSteps(t *testing.T) {
	const input = `10R5L5R10L4R5L5`

	var expect = []Step{
		{isMove: true, distance: 10},
		{rotateClockwise: true},
		{isMove: true, distance: 5},
		{rotateClockwise: false},
		{isMove: true, distance: 5},
		{rotateClockwise: true},
		{isMove: true, distance: 10},
		{rotateClockwise: false},
		{isMove: true, distance: 4},
		{rotateClockwise: true},
		{isMove: true, distance: 5},
		{rotateClockwise: false},
		{isMove: true, distance: 5},
	}

	steps, err := parseSteps([]byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(steps), 13; got != want {
		t.Errorf("len(steps) = %d, want %d", got, want)
	}

	if !reflect.DeepEqual(steps, expect) {
		t.Errorf("steps = %v, want %v", steps, expect)
	}

}

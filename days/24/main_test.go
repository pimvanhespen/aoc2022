package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const testInput = `#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`

var r = strings.NewReplacer(">", "@",
	"<", "@",
	"^", "@",
	"v", "@",
	"2", "@",
	"3", "@",
	"4", "@",
	"5", "@",
	"6", "@",
	"7", "@",
	"8", "@",
	"9", "@")

func TestParse(t *testing.T) {
	valley, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	expect := r.Replace(testInput)

	got := valley.StringWhen(0)
	if got != expect {
		t.Errorf("Expected\n%s\ngot\n%s\n", testInput, got)
	}
}

func TestValley_Next(t *testing.T) {
	begin := `#E######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`
	valley, err := parse(strings.NewReader(begin))
	if err != nil {
		t.Fatal(err)
	}

	rows := bytes.Split([]byte(example), []byte{'\n'})

	for i := 0; i < len(rows); i += 8 {
		got := valley.StringWhen(i / 8) // ignore the initial state

		want := string(bytes.Join(rows[i+1:i+7], []byte{'\n'}))
		want = strings.Replace(want, "E", ".", -1)
		want = r.Replace(want)

		if got != want {
			t.Errorf("Expected\n%s\ngot\n%s\n", want, got)
		}
	}
}

func TestSolve1(t *testing.T) {
	const want = 18
	valley, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	actual := solve1(valley)
	if actual != want {
		t.Errorf("Expected %d, got %d", want, actual)
	}
}

var example = `Initial state:
#E######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#

Minute 1, move down:
#.######
#E>3.<.#
#<..<<.#
#>2.22.#
#>v..^<#
######.#

Minute 2, move down:
#.######
#.2>2..#
#E^22^<#
#.>2.^>#
#.>..<.#
######.#

Minute 3, wait:
#.######
#<^<22.#
#E2<.2.#
#><2>..#
#..><..#
######.#

Minute 4, move up:
#.######
#E<..22#
#<<.<..#
#<2.>>.#
#.^22^.#
######.#

Minute 5, move right:
#.######
#2Ev.<>#
#<.<..<#
#.^>^22#
#.2..2.#
######.#

Minute 6, move right:
#.######
#>2E<.<#
#.2v^2<#
#>..>2>#
#<....>#
######.#

Minute 7, move down:
#.######
#.22^2.#
#<vE<2.#
#>>v<>.#
#>....<#
######.#

Minute 8, move left:
#.######
#.<>2^.#
#.E<<.<#
#.22..>#
#.2v^2.#
######.#

Minute 9, move up:
#.######
#<E2>>.#
#.<<.<.#
#>2>2^.#
#.v><^.#
######.#

Minute 10, move right:
#.######
#.2E.>2#
#<2v2^.#
#<>.>2.#
#..<>..#
######.#

Minute 11, wait:
#.######
#2^E^2>#
#<v<.^<#
#..2.>2#
#.<..>.#
######.#

Minute 12, move down:
#.######
#>>.<^<#
#.<E.<<#
#>v.><>#
#<^v^^>#
######.#

Minute 13, move down:
#.######
#.>3.<.#
#<..<<.#
#>2E22.#
#>v..^<#
######.#

Minute 14, move right:
#.######
#.2>2..#
#.^22^<#
#.>2E^>#
#.>..<.#
######.#

Minute 15, move right:
#.######
#<^<22.#
#.2<.2.#
#><2>E.#
#..><..#
######.#

Minute 16, move right:
#.######
#.<..22#
#<<.<..#
#<2.>>E#
#.^22^.#
######.#

Minute 17, move down:
#.######
#2.v.<>#
#<.<..<#
#.^>^22#
#.2..2E#
######.#

Minute 18, move down:
#.######
#>2.<.<#
#.2v^2<#
#>..>2>#
#<....>#
######E#`

func TestValley_IsValid(t *testing.T) {
	const data = `
#####
#...#
#####`

	valley, err := parse(strings.NewReader(data[1:]))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(valley.StringWhen(0))

	invalidVectors := []Vector2D{
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: -1},
		{X: 2, Y: -1},
		{X: 3, Y: -1},
		{X: 3, Y: 0},
		{X: 3, Y: 1},
	}

	for _, v := range invalidVectors {
		if valley.IsValid(State{
			Position: v,
		}) {
			t.Errorf("Expected %v to be invalid", v)
		}
	}

}

package main

import (
	"fmt"
	"strings"
	"testing"
)

const testData = `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`

func TestSolve1(t *testing.T) {
	v, err := parse(strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}

	const want = 1651
	got := solve1(v)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	v, err := parse(strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}

	const want = 1707
	got := solve2(v)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSolveBounds(t *testing.T) {

	c := &Valve{
		Name:     "C",
		FlowRate: 2,
	}

	b := &Valve{
		Name:     "B",
		FlowRate: 5,
		Dist: map[*Valve]int{
			c: 1,
		},
	}

	a := &Valve{
		Name:     "A",
		FlowRate: 0,
		Dist: map[*Valve]int{
			b: 29,
		},
	}

	const want = 0
	got := solve1(a)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func PrintValves(v *Valve, indent string, m map[*Valve]struct{}) {
	if m == nil {
		m = map[*Valve]struct{}{}
	}
	fmt.Printf("%s%s\n", indent, v)

	if _, ok := m[v]; ok {
		return
	}
	m[v] = struct{}{}

	for _, t := range v.Tunnels {
		PrintValves(t, indent+"  ", m)
	}
}

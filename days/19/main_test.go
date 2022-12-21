package main

import (
	"reflect"
	"strings"
	"testing"
)

var testInput = `

Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`

func TestSolve1(t *testing.T) {
	bps, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range bps {
		t.Log(b)
	}

	if got, want := solve1(bps), 33; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestBlueprint_Simulate(t *testing.T) {
	bps, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := bps[0].Simulate(32), 56; got.Stock.Geode != want {
		t.Errorf("got %d, want %d", got.Stock.Geode, want)
	}

}
func TestBlueprint_Simulate2(t *testing.T) {
	bps, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	initialState := State{
		Turns:  32,
		Stock:  Resources{},
		Robots: Resources{Ore: 1},
	}
	if got, want := bps[1].simulate(initialState, 47), 62; got.Stock.Geode != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestSolve2(t *testing.T) {
	bps, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range bps {
		t.Log(b)
	}

	if len(bps) > 2 {
		bps = bps[:3]
	}

	if got, want := solve2(bps), 62*56; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestState_Wait(t *testing.T) {
	type fields struct {
		Turns  int
		Stock  Resources
		Robots Resources
	}
	tests := []struct {
		name   string
		fields fields
		want   State
	}{
		{
			name: "wait",
			fields: fields{
				Turns:  1,
				Stock:  Resources{Ore: 0},
				Robots: Resources{Ore: 1},
			},
			want: State{
				Turns:  0,
				Stock:  Resources{Ore: 1},
				Robots: Resources{Ore: 1},
			},
		}, {
			name: "wait 2 turns",
			fields: fields{
				Turns:  2,
				Stock:  Resources{Ore: 0},
				Robots: Resources{Ore: 5},
			},
			want: State{
				Turns:  0,
				Stock:  Resources{Ore: 10},
				Robots: Resources{Ore: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{
				Turns:  tt.fields.Turns,
				Stock:  tt.fields.Stock,
				Robots: tt.fields.Robots,
			}
			if got := s.Wait(s.Turns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wait() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Next(t *testing.T) {
	type fields struct {
		Turns  int
		Stock  Resources
		Robots Resources
	}
	type args struct {
		cost     Resources
		addRobot Resources
		wait     int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   State
	}{
		{
			name: "build 1 ore robot",
			fields: fields{
				Turns:  1,
				Stock:  Resources{Ore: 4},
				Robots: Resources{Ore: 1},
			},
			args: args{
				cost:     Resources{Ore: 4},
				addRobot: Resources{Ore: 1},
				wait:     0,
			},
			want: State{
				Turns:  0,
				Stock:  Resources{Ore: 1},
				Robots: Resources{Ore: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{
				Turns:  tt.fields.Turns,
				Stock:  tt.fields.Stock,
				Robots: tt.fields.Robots,
			}
			if got := s.buyAfter(tt.args.wait, tt.args.cost, tt.args.addRobot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buy() = %v, want %v", got, tt.want)
			}
		})
	}
}

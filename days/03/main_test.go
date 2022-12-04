package main

import (
	"strings"
	"testing"
)

var example = `
vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw`

func TestSolve1(t *testing.T) {
	rucksacks, err := parseInput(strings.NewReader(example))
	if err != nil {
		t.Fatal(err)
	}

	result, err := solve1(rucksacks)
	if err != nil {
		t.Fatal(err)
	}

	if result != 157 {
		t.Fatalf("expected 157, got %d", result)
	}
}

func Test_value(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "a",
			args: args{b: 'a'},
			want: 1,
		},
		{
			name: "A",
			args: args{b: 'A'},
			want: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := value(tt.args.b); got != tt.want {
				t.Errorf("value() = %v, want %v", got, tt.want)
			}
		})
	}
}

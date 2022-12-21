package main

import (
	"strings"
	"testing"
)

var testInput = `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

func TestSolve1(t *testing.T) {
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	const want = 152

	res := solve1(input)
	if res != want {
		t.Errorf("solve1() = %v, want %v", res, want)
	}
}

func TestSolve2(t *testing.T) {
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}

	const want = 301

	res := solve2(input)
	if res != want {
		t.Errorf("solve2() = %v, want %v", res, want)
	}
}

func Test_reverse(t *testing.T) {
	type args struct {
		sum    int
		op     rune
		value  int
		isLeft bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "add",
			args: args{
				sum:    10,
				op:     '+',
				value:  9,
				isLeft: true,
			},
			want: 1,
		},
		{
			name: "sub",
			args: args{
				sum:    10,
				op:     '-',
				value:  5,
				isLeft: true,
			},
			want: 15,
		},
		{
			name: "mul",
			args: args{
				sum:    10,
				op:     '*',
				value:  5,
				isLeft: true,
			},
			want: 2,
		},
		{
			name: "div",
			args: args{
				sum:    10,
				op:     '/',
				value:  5,
				isLeft: true,
			},
			want: 50,
		},
		{
			name: "add",
			args: args{
				sum:    10,
				op:     '+',
				value:  6,
				isLeft: false,
			},
			want: 4,
		},
		{
			name: "sub",
			args: args{
				sum:    10,
				op:     '-',
				value:  5,
				isLeft: false,
			},
			want: -5,
		},
		{
			name: "mul",
			args: args{
				sum:    10,
				op:     '*',
				value:  5,
				isLeft: false,
			},
			want: 2,
		},
		{
			name: "div",
			args: args{
				sum:    10,
				op:     '/',
				value:  50,
				isLeft: false,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.args.sum, tt.args.op, tt.args.value, tt.args.isLeft); got != tt.want {
				t.Errorf("reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

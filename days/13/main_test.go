package main

import (
	"strings"
	"testing"
)

const testInput = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

func TestSolve1(t *testing.T) {
	packets := parse(strings.NewReader(testInput))
	num := solve1(packets)

	if num != 13 {
		t.Errorf("solve1() = %v, want %v", num, 13)
	}
}

func TestSolve2(t *testing.T) {
	packets := parse(strings.NewReader(testInput))
	num := solve2(packets)

	if num != 140 {
		t.Errorf("solve2() = %v, want %v", num, 140)
	}
}

func Test_parsePacket(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{
			name: "empty",
			arg:  "[]",
		},
		{
			name: "single",
			arg:  "[1]",
		},
		{
			name: "multiple",
			arg:  "[1,2,3]",
		},
		{
			name: "nested",
			arg:  "[1,[2,3]]",
		},
		{
			name: "nested multiple",
			arg:  "[1,[2,3],[4,5]]",
		},
		{
			name: "big one",
			arg:  "[[[[8,6,7],9,7,10],[[2,2,4],0,[4,9,10],[4,8,1,1],9],5],[8,8,[5,7],1,3],[],[6,[1,[0,1],[6,10,9]]],[]]",
		},
		{
			name: "big one",
			arg:  "[[],[],[4,[8],[4,[4,4,8,10,6],8,9],5]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := parsePacket(tt.arg)
			if got.String() != tt.arg {
				t.Errorf("parseListItem() got = %v, want %v", got, tt.arg)
			}
		})
	}
}

func Test_compareItem(t *testing.T) {
	type args struct {
		a Value
		b Value
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty - equal",
			args: args{
				a: ListValue{},
				b: ListValue{},
			},
			want: 0,
		},
		{
			name: "single - equal",
			args: args{
				a: ListValue{IntValue(1)},
				b: ListValue{IntValue(1)},
			},
			want: 0,
		},
		{
			name: "multiple - equal",
			args: args{
				a: ListValue{IntValue(1), IntValue(2), IntValue(3)},
				b: ListValue{IntValue(1), IntValue(2), IntValue(3)},
			},
			want: 0,
		},
		{
			name: "nested - equal",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}},
			},
			want: 0,
		},
		{
			name: "nested multiple - equal",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(5)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(5)}},
			},
			want: 0,
		},
		{
			name: "empty - not equal",
			args: args{
				a: ListValue{},
				b: ListValue{IntValue(1)},
			},
			want: -1,
		},
		{
			name: "single - not equal",
			args: args{
				a: ListValue{IntValue(1)},
				b: ListValue{IntValue(2)},
			},
			want: -1,
		},
		{
			name: "multiple - not equal",
			args: args{
				a: ListValue{IntValue(1), IntValue(2), IntValue(3)},
				b: ListValue{IntValue(1), IntValue(2), IntValue(4)},
			},
			want: -1,
		},
		{
			name: "nested - not equal",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(4)}},
			},
			want: -1,
		},
		{
			name: "nested multiple - not equal",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(5)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(6)}},
			},
			want: -1,
		},
		{
			name: "empty - more",
			args: args{
				a: ListValue{IntValue(1)},
				b: ListValue{},
			},
			want: 1,
		},
		{
			name: "single - more",
			args: args{
				a: ListValue{IntValue(2)},
				b: ListValue{IntValue(1)},
			},
			want: 1,
		},
		{
			name: "multiple - more",
			args: args{
				a: ListValue{IntValue(1), IntValue(2), IntValue(4)},
				b: ListValue{IntValue(1), IntValue(2), IntValue(3)},
			},
			want: 1,
		},
		{
			name: "nested - more",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(4)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}},
			},
			want: 1,
		},
		{
			name: "nested multiple - more",
			args: args{
				a: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(6)}},
				b: ListValue{IntValue(1), ListValue{IntValue(2), IntValue(3)}, ListValue{IntValue(4), IntValue(5)}},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.a.Compare(tt.args.b); got != tt.want {
				t.Errorf("compareItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

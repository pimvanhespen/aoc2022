package main

import (
	"reflect"
	"strings"
	"testing"
)

const testInput = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestSolve1(t *testing.T) {
	input, err := parse(strings.NewReader(testInput))
	if err != nil {
		t.Fatal(err)
	}
	got := solve1(input, 2022)
	want := 3068
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestBlock_Point(t *testing.T) {
	type fields struct {
		Position Vec2
		Shape    Shape
	}
	tests := []struct {
		name   string
		fields fields
		want   []Vec2
	}{
		{
			name: "test",
			fields: fields{
				Position: Vec2{X: 1, Y: 2},
				Shape: Shape{
					fields: []Vec2{
						{X: 0, Y: 0},
						{X: 1, Y: 0},
						{X: 0, Y: 1},
					},
				},
			},
			want: []Vec2{
				{X: 1, Y: 2},
				{X: 2, Y: 2},
				{X: 1, Y: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Block{
				Position: tt.fields.Position,
				Shape:    tt.fields.Shape,
			}
			if got := b.Point(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_String(t *testing.T) {
	const want = `|....#..|
|...###.|
|....#..|
|...####|
+-------+`
	field := Field{width: 7}
	field.Add(HorizontalLine, Vec2{3, 0})
	field.Add(Plus, Vec2{3, 1})

	s := field.String()
	if s != want {
		var sb strings.Builder
		sb.WriteString("got:\t\twant:\n")
		pw, pg := strings.Split(want, "\n"), strings.Split(s, "\n")
		for i := 0; i < max(len(pw), len(pg)); i++ {
			if i < len(pg) {
				sb.WriteString(pg[i])
			} else {
				sb.WriteString(strings.Repeat(" ", len(pw[i])))
			}
			sb.WriteByte('\t')
			if i < len(pw) {
				sb.WriteString(pw[i])
			}
			sb.WriteByte('\n')
		}
		t.Errorf(sb.String())
	}
}

func TestBounds(t *testing.T) {
	type args struct {
		block Block
	}

	tests := []struct {
		name string
		args args
		want Bounds
	}{
		{
			name: "test",
			args: args{
				block: Block{Position: Vec2{X: 0, Y: 0}, Shape: Plus},
			},
			want: Bounds{minX: 0, maxX: 2, minY: 0, maxY: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.block.Bounds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

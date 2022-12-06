package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func Test_solve1(t *testing.T) {
	type args struct {
		bts []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example1",
			args: args{
				bts: []byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`),
			},
			want: 5,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`nppdvjthqldpwncqszvftbrmjlhg`),
			},
			want: 6,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`),
			},
			want: 10,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`),
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve1(tt.args.bts); got != tt.want {
				t.Errorf("solve1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solve_map(t *testing.T) {
	type args struct {
		bts []byte
		len int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example1",
			args: args{
				bts: []byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`),
				len: 4,
			},
			want: 5,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`nppdvjthqldpwncqszvftbrmjlhg`),
				len: 4,
			},
			want: 6,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`),
				len: 4,
			},
			want: 10,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`),
				len: 4,
			},
			want: 11,
		},

		{
			name: "example2",
			args: args{
				bts: []byte(`mjqjpqmgbljsphdztnvjfqwrcgsmlb`),
				len: 14,
			},
			want: 19,
		}, {
			name: "example2",
			args: args{
				bts: []byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`),
				len: 14,
			},
			want: 23,
		}, {
			name: "example2",
			args: args{
				bts: []byte(`nppdvjthqldpwncqszvftbrmjlhg`),
				len: 14,
			},
			want: 23,
		}, {
			name: "example2",
			args: args{
				bts: []byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`),
				len: 14,
			},
			want: 29,
		},
		{
			name: "example2",
			args: args{
				bts: []byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`),
				len: 14,
			},
			want: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve_map(tt.args.bts, tt.args.len); got != tt.want {
				t.Errorf("solve1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solve2(t *testing.T) {
	type args struct {
		bts []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example1",
			args: args{
				bts: []byte(`mjqjpqmgbljsphdztnvjfqwrcgsmlb`),
			},
			want: 19,
		}, {
			name: "example1",
			args: args{
				bts: []byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`),
			},
			want: 23,
		}, {
			name: "example1",
			args: args{
				bts: []byte(`nppdvjthqldpwncqszvftbrmjlhg`),
			},
			want: 23,
		}, {
			name: "example1",
			args: args{
				bts: []byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`),
			},
			want: 29,
		},
		{
			name: "example1",
			args: args{
				bts: []byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`),
			},
			want: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve2(tt.args.bts); got != tt.want {
				t.Errorf("solve2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSolve1(b *testing.B) {

	f, err := os.Open("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	input, err := io.ReadAll(f)
	if err != nil {
		b.Error(err)
	}

	var total int

	b.Run("allDifferent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			total += solve1(input)
		}
	})

	b.Run("map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			total += solve_map(input, 4)
		}
	})

	fmt.Println(total)
}

func BenchmarkSolve2(b *testing.B) {

	f, err := os.Open("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	input, err := io.ReadAll(f)
	if err != nil {
		b.Error(err)
	}

	var total int

	b.Run("allDifferent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			total += solve2(input)
		}
	})

	b.Run("map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			total += solve_map(input, 14)
		}
	})

	fmt.Println(total)
}

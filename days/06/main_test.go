package main

import "testing"

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

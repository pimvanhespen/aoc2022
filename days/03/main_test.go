package main

import (
	"bytes"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/list"
	"io"
	"os"
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

func BenchmarkSolve2(b *testing.B) {

	f, err := os.Open("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			b.Fatal(e)
		}
	}()

	bts, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}

	bs, err := parseInput(bytes.NewReader(bts))
	if err != nil {
		b.Fatal(err)
	}

	var total int
	var n int

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("initial", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2(bs)
			total += n
		}
	})

	b.Run("speed-up", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_speed(bts)
			total += n
		}
	})

	b.Run("mem-up", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_mem(bts)
			total += n
		}
	})

	fmt.Println(total)
}

func BenchmarkSolve1(b *testing.B) {

	f, err := os.Open("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	bs, err := parseInput(f)
	if err != nil {
		b.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		b.Fatal(err)
	}

	ss := list.Transform(bs, func(b []byte) string {
		return string(b)
	})

	line := strings.Join(ss, "\n") + "\n"

	bts := []byte(line)

	var total int
	var n int

	b.ResetTimer()

	b.Run("pim", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1(bs)
			total += n
		}
	})

	b.Run("lowMem", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1_lowMem(bts)
			total += n
		}
	})

	b.Run("indexByte", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1_lowMem_IndexByte(bts)
			total += n
		}
	})

	b.Run("structMap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1_lowMem_structMap(bts)
			total += n
		}
	})

	b.Run("speed", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1_speed(bts)
			total += n
		}
	})

	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve1_stdlib(bts)
			total += n
		}
	})

	fmt.Println(total)
}

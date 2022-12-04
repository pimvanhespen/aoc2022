package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

// BenchmarkSolve2
// goos: darwin
// goarch: arm64
// pkg: github.com/pimvanhespen/aoc2022/days/03
// BenchmarkSolve2/solve2
// BenchmarkSolve2/solve2-8         	    1812	    655702 ns/op	   94708 B/op	    2117 allocs/op
// BenchmarkSolve2/solve2_speed
// BenchmarkSolve2/solve2_speed-8   	    3310	    362484 ns/op	   29967 B/op	     712 allocs/op
// BenchmarkSolve2/solve2_mem
// BenchmarkSolve2/solve2_mem-8     	    4095	    298299 ns/op	     662 B/op	      27 allocs/op
// BenchmarkSolve2/solve2_bytes
// BenchmarkSolve2/solve2_bytes-8   	   77505	     15466 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSolve2/solve2_bytes_ptr
// BenchmarkSolve2/solve2_bytes_ptr-8      82242	     14568 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSolve2/solve2_bytes_ref
// BenchmarkSolve2/solve2_bytes_ref-8      81895	     14618 ns/op	       0 B/op	       0 allocs/op
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

	b.Run("solve2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2(bs)
			total += n
		}
	})

	b.Run("solve2_speed", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_speed(bts)
			total += n
		}
	})

	b.Run("solve2_mem", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_mem(bts)
			total += n
		}
	})

	b.Run("solve2_bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_bytes(bts)
			total += n
		}
	})

	b.Run("solve2_bytes_ptr", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_bytes_ptr(bts)
			total += n
		}
	})

	b.Run("solve2_bytes_ref", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			n, _ = solve2_bytes_ref(bts)
			total += n
		}
	})

	fmt.Println(total)
}

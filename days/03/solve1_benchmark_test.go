package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/list"
	"os"
	"strings"
	"testing"
)

// BenchmarkSolve1
// goos: darwin
// goarch: arm64
// pkg: github.com/pimvanhespen/aoc2022/days/03
// BenchmarkSolve1/pim
// BenchmarkSolve1/pim-8         	    1699	    700079 ns/op	  132519 B/op	    3902 allocs/op
// BenchmarkSolve1/lowMem
// BenchmarkSolve1/lowMem-8      	    1137	   1052060 ns/op	    1225 B/op	      20 allocs/op
// BenchmarkSolve1/indexByte
// BenchmarkSolve1/indexByte-8   	    1137	   1051328 ns/op	    1224 B/op	      20 allocs/op
// BenchmarkSolve1/structMap
// BenchmarkSolve1/structMap-8   	    3703	    322972 ns/op	     251 B/op	       9 allocs/op
// BenchmarkSolve1/speed
// BenchmarkSolve1/speed-8       	    4226	    283357 ns/op	    5599 B/op	     121 allocs/op
// BenchmarkSolve1/bytes
// BenchmarkSolve1/bytes-8       	   60631	     19754 ns/op	       0 B/op	       0 allocs/op
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
			n, _ = solve1_noAllocs(bts)
			total += n
		}
	})

	fmt.Println(total)
}

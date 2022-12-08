package main

import (
	"bytes"
	"os"
	"testing"
)

func Test_solve1(t *testing.T) {
	const input = `30373
25512
65332
33549
35390
`
	const want = 21

	f, err := parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}

	got, err := solve1(f)
	if err != nil {
		t.Fatal(err)
	}

	if isVisible(f, 2, 2) {
		t.Errorf("expected (2, 2) to not be visible")
	}
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func BenchmarkSolve1(b *testing.B) {
	//	const input = `30373
	//25512
	//65332
	//33549
	//35390
	//`

	fl, err := os.Open("input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer fl.Close()

	f, err := parse(fl)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("solve1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			t, sErr := solve1(f)
			if sErr != nil {
				b.Fatal(sErr)
			}
			if t != 1647 {
				b.Errorf("got %d, want %d", t, 21)
			}
		}
	})

	b.Run("solve1_visibility_matrix", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			t, sErr := solve1_visibility_matrix(f)
			if sErr != nil {
				b.Fatal(sErr)
			}
			if t != 1647 {
				b.Errorf("got %d, want %d", t, 21)
			}
		}
	})
}

func Test_solve1_fast(t *testing.T) {
	const input = `30373
25512
65332
33549
35390
`
	const want = 21

	f, err := parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}

	got, err := solve1_visibility_matrix(f)
	if err != nil {
		t.Fatal(err)
	}

	if isVisible(f, 2, 2) {
		t.Errorf("expected (2, 2) to not be visible")
	}
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func Test_solve2(t *testing.T) {
	const input = `30373
25512
65332
33549
35390
`
	const want = 8

	f, err := parse(bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}

	got, err := solve2(f)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if 8 != scenicScore(f, 2, 3) {
		t.Errorf("expected (2, 3) to have a score of 8, got %d", scenicScore(f, 3, 2))
	}

	up := scenicScoreInDirection(f, 2, 3, 0, -1)
	if 2 != up {
		t.Errorf("expected up (2, 3) to have a score of 2, got %d", up)
	}

	down := scenicScoreInDirection(f, 2, 3, 0, 1)
	if 1 != down {
		t.Errorf("expected down (2, 3) to have a score of 1, got %d", down)
	}

	left := scenicScoreInDirection(f, 2, 3, -1, 0)
	if 2 != left {
		t.Errorf("expected left (2, 3) to have a score of 2, got %d", left)
	}

	right := scenicScoreInDirection(f, 2, 3, 1, 0)
	if 2 != right {
		t.Errorf("expected right (2, 3) to have a score of 2, got %d", right)
	}
}

func TestIsVisbile(t *testing.T) {
	grid := [][]byte{
		[]byte("222"),
		[]byte("212"),
		[]byte("222"),
	}

	if isVisible(grid, 1, 1) {
		t.Errorf("expected (1, 1) to not be visible")
	}

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			if x == 1 && y == 1 {
				continue
			}
			if !isVisible(grid, x, y) {
				t.Errorf("expected (%d, %d) to be visible", x, y)
			}
		}
	}
}

package main

import (
	"aoc2022/internal/puzzleinput"
	"bufio"
	"fmt"
	"io"
)

type Hand int

const (
	Rock = iota + 1
	Paper
	Scissors
)

type Outcome int

const (
	Loss = 0
	Draw = 3
	Win  = 6
)

func outcomeFromRune(r rune) Outcome {
	switch r {
	case 'X':
		return Loss
	case 'Y':
		return Draw
	case 'Z':
		return Win
	default:
		panic("unknown outcome")
	}
}

// Outcome is the outcome for this hand compared to the other hand.
// 6 on win
// 3 on draw
// 0 on loss
func (h Hand) Outcome(other Hand) Outcome {
	if h == other {
		return Draw
	}

	if h > other && other != Rock {
		return Win
	}

	if h == Rock && other == Scissors {
		return Win
	}

	if h == Paper && other == Rock {
		return Win
	}

	return Loss
}

func (h Hand) CalcHandForOutcome(outcome Outcome) Hand {
	switch outcome {
	case Loss:
		switch h {
		case Rock:
			return Scissors
		case Paper:
			return Rock
		case Scissors:
			return Paper
		}
	case Draw:
		return h

	case Win:
		switch h {
		case Rock:
			return Paper
		case Paper:
			return Scissors
		case Scissors:
			return Rock
		}
	}
	panic("unknown outcome")
}

type Round struct {
	Left, Right rune
}

func main() {
	input, err := puzzleinput.Get(2)
	if err != nil {
		panic(err)
	}

	invs, err := parseInput(input)
	if err != nil {
		panic(err)
	}

	result1, err := solve1(invs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", result1)

	result2, err := solve2(invs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result2)

}

func parseInput(rc io.ReadCloser) ([]Round, error) {
	defer rc.Close()

	var rounds []Round

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var left, right rune

		_, err := fmt.Sscanf(line, "%c %c", &left, &right)
		if err != nil {
			return nil, err
		}

		rounds = append(rounds, Round{
			left, right,
		})

	}

	return rounds, nil
}

func handFromRune(left rune) Hand {
	switch left {
	case 'A', 'X':
		return Rock
	case 'B', 'Y':
		return Paper
	case 'C', 'Z':
		return Scissors
	default:
		panic("unknown hand")
	}
}

func solve1(rounds []Round) (int, error) {
	var total int
	for _, round := range rounds {
		elf, self := handFromRune(round.Left), handFromRune(round.Right)
		total += int(self) + int(self.Outcome(elf))
	}

	return total, nil
}

func solve2(rounds []Round) (int, error) {
	var total int
	for _, round := range rounds {
		result, err := calcRoundScore_Strategy2(round)
		if err != nil {
			return 0, err
		}
		total += result
	}
	return total, nil
}

func calcRoundScore_Strategy2(r Round) (int, error) {
	elf := handFromRune(r.Left)
	outcome := outcomeFromRune(r.Right)
	self := elf.CalcHandForOutcome(outcome)

	return int(self) + int(outcome), nil
}

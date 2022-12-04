package main

import (
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/list"
	"github.com/pimvanhespen/aoc2022/pkg/rps"
)

type inputRow struct {
	left, right rune
}

func main() {

	in, err := aoc.Get(2)
	if err != nil {
		panic(err)
	}

	parser := aoc.Parser[inputRow]{
		SkipEmptyLines: true,
		ParseFn:        parseLine,
	}

	invs, err := parser.Rows(in)
	if err != nil {
		panic(err)
	}

	result1, err := solve(invs, transformA, strategyA)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", result1)

	result2, err := solve(invs, transformB, strategyB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result2)

}

func parseLine(line string) (inputRow, error) {
	var input inputRow

	_, err := fmt.Sscanf(line, "%c %c", &input.left, &input.right)
	if err != nil {
		return inputRow{}, err
	}

	return input, nil
}

func parseOutcome(r rune) (rps.Outcome, error) {
	switch r {
	case 'X':
		return rps.Loss, nil
	case 'Y':
		return rps.Draw, nil
	case 'Z':
		return rps.Win, nil
	default:
		return rps.Outcome{}, fmt.Errorf("unknown outcome '%c'", r)
	}
}

func parseHand(r rune) (rps.Hand, error) {
	switch r {
	case 'A', 'X':
		return rps.Rock, nil
	case 'B', 'Y':
		return rps.Paper, nil
	case 'C', 'Z':
		return rps.Scissors, nil
	default:
		return rps.Hand{}, fmt.Errorf("unknown hand '%c'", r)
	}
}

func outcomeValue(outcome rps.Outcome) int {
	switch outcome {
	case rps.Loss:
		return 0
	case rps.Draw:
		return 3
	case rps.Win:
		return 6
	default:
		panic("unknown outcome")
	}
}

func handValue(hand rps.Hand) int {
	switch hand {
	case rps.Rock:
		return 1
	case rps.Paper:
		return 2
	case rps.Scissors:
		return 3
	default:
		panic("unknown hand")
	}
}

func solve[T any](
	rows []inputRow,
	rowToInput func(inputRow) (T, error),
	strategy func(s T) (int, error)) (int, error) {

	rounds, err := list.TransformErr(rows, rowToInput)
	if err != nil {
		return 0, err
	}

	var total int
	for _, round := range rounds {
		result, solveErr := strategy(round)
		if solveErr != nil {
			return 0, solveErr
		}
		total += result
	}

	return total, nil
}

// -- solve 1 --

type inputA struct {
	OpponentHand rps.Hand
	PlayerHand   rps.Hand
}

func transformA(r inputRow) (inputA, error) {
	opponentsHand, err := parseHand(r.left)
	if err != nil {
		return inputA{}, err
	}
	playerHand, err := parseHand(r.right)
	if err != nil {
		return inputA{}, err
	}

	input := inputA{
		OpponentHand: opponentsHand,
		PlayerHand:   playerHand,
	}
	return input, nil
}

func strategyA(i inputA) (int, error) {

	outcome, err := i.PlayerHand.Play(i.OpponentHand)
	if err != nil {
		return 0, err
	}

	res := outcomeValue(outcome) + handValue(i.PlayerHand)
	return res, nil
}

type inputB struct {
	OpponentHand   rps.Hand
	DesiredOutcome rps.Outcome
}

func transformB(r inputRow) (inputB, error) {
	opponentsHand, err := parseHand(r.left)
	if err != nil {
		return inputB{}, err
	}

	desiredOutcome, err := parseOutcome(r.right)
	if err != nil {
		return inputB{}, err
	}

	input := inputB{
		OpponentHand:   opponentsHand,
		DesiredOutcome: desiredOutcome,
	}
	return input, nil
}

func strategyB(i inputB) (int, error) {
	handToPlay, err := rps.HandForOutcome(i.DesiredOutcome, i.OpponentHand)
	if err != nil {
		return 0, err
	}

	res := outcomeValue(i.DesiredOutcome) + handValue(handToPlay)
	return res, nil
}

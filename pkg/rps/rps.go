// Package rps implements the logic for a simple rock-paper-scissors game.
package rps

import "fmt"

var (
	ErrInvalidHand    = newRockPaperScissorsError("invalid hand")
	ErrInvalidOutcome = newRockPaperScissorsError("invalid outcome")
)

var (
	Loss = Outcome{"Loss"}
	Draw = Outcome{"Draw"}
	Win  = Outcome{"Win"}
)

type Outcome struct {
	name string
}

func (o Outcome) String() string {
	return o.name
}

func (o Outcome) isZero() bool {
	return o.name == ""
}

var (
	Rock     = Hand{rule: rockRule}
	Paper    = Hand{rule: paperRule}
	Scissors = Hand{rule: scissorsRule}
)

type Hand struct {
	rule *rule
}

// Play calculates the outcome of a game of rock-paper-scissors.
// It returns an error if either hand is invalid.
func (h Hand) Play(opponent Hand) (Outcome, error) {
	if h.isInvalid() || opponent.isInvalid() {
		return Outcome{}, ErrInvalidHand
	}

	return h.rule.outcome(opponent.rule), nil
}

func (h Hand) String() string {
	return h.rule.name
}

func (h Hand) isInvalid() bool {
	return h.rule == nil
}

// HandForOutcome calculates which hand should be played against the opponents hand
// to get the desired outcome.
func HandForOutcome(outcome Outcome, opponent Hand) (Hand, error) {

	if opponent.isInvalid() {
		return Hand{}, ErrInvalidHand
	}

	switch outcome {
	case Loss:
		return Hand{opponent.rule.win}, nil
	case Draw:
		return Hand{opponent.rule}, nil
	case Win:
		return Hand{opponent.rule.loss}, nil
	default:
		return Hand{}, ErrInvalidOutcome
	}
}

// -- Internal --

func init() {
	initRules()
}

func initRules() {
	rockRule.win = scissorsRule
	rockRule.loss = paperRule

	paperRule.win = rockRule
	paperRule.loss = scissorsRule

	scissorsRule.win = paperRule
	scissorsRule.loss = rockRule
}

var (
	rockRule     = newRule("Rock")
	paperRule    = newRule("Paper")
	scissorsRule = newRule("Scissors")
)

type rule struct {
	name string
	win  *rule
	loss *rule
}

func newRule(name string) *rule {
	return &rule{name: name}
}

func (h *rule) outcome(other *rule) Outcome {
	switch other {
	case h.win:
		return Win
	case h:
		return Draw
	case h.loss:
		return Loss
	default:
		panic(fmt.Sprintf("unknown rule %v", other))
	}
}

type RockPaperScissorsError struct {
	msg string
}

func newRockPaperScissorsError(msg string) *RockPaperScissorsError {
	return &RockPaperScissorsError{msg: msg}
}

func (e RockPaperScissorsError) Error() string {
	return e.msg
}

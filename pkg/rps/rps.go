// Package rps implements the logic for a simple rock-paper-scissors game.
package rps

import "fmt"

var (
	Loss = Outcome{"Loss"}
	Draw = Outcome{"Draw"}
	Win  = Outcome{"Win"}
)

type Outcome struct {
	name string
}

var (
	Rock     = Hand{rule: rockRule}
	Paper    = Hand{rule: paperRule}
	Scissors = Hand{rule: scissorsRule}
)

type Hand struct {
	rule *rule
}

func (h Hand) Outcome(other Hand) (Outcome, error) {
	return h.rule.outcome(other.rule)
}

// HandForOutcome calculates which hand should be played against the oppnents hand
// to get the desired outcome.
func HandForOutcome(outcome Outcome, opponent Hand) (Hand, error) {
	switch outcome {
	case Loss:
		return ptrToHand(opponent.rule.win)
	case Draw:
		return ptrToHand(opponent.rule)
	case Win:
		return ptrToHand(opponent.rule.loss)
	}
	return Hand{}, fmt.Errorf("unknown outcome %v", outcome)
}

// -- Internal --

func init() {
	initHands()
}

func initHands() {
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

// Outcome is the outcome for this rule compared to the other rule.
// 6 on win, 3 on draw, 0 on loss.
func (h *rule) outcome(other *rule) (Outcome, error) {
	switch other {
	case h.win:
		return Win, nil
	case h:
		return Draw, nil
	case h.loss:
		return Loss, nil
	}

	return Outcome{}, fmt.Errorf("unknown outcome for %v vs %v", h, other)
}

func ptrToHand(h *rule) (Hand, error) {
	switch h {
	case rockRule:
		return Rock, nil
	case paperRule:
		return Paper, nil
	case scissorsRule:
		return Scissors, nil
	}
	return Hand{}, fmt.Errorf("unknown rule %v", h)
}

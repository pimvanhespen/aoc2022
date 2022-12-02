package rps

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initHands()

	code := m.Run()

	os.Exit(code)
}

func TestHand_Outcome(t *testing.T) {
	type args struct {
		other Hand
	}
	tests := []struct {
		name string
		h    Hand
		args args
		want Outcome
	}{
		{
			name: "Rock vs Rock",
			h:    Rock,
			args: args{other: Rock},
			want: Draw,
		},
		{
			name: "Rock vs Paper",
			h:    Rock,
			args: args{other: Paper},
			want: Loss,
		},
		{
			name: "Rock vs Scissors",
			h:    Rock,
			args: args{other: Scissors},
			want: Win,
		},
		{
			name: "Paper vs Rock",
			h:    Paper,
			args: args{other: Rock},
			want: Win,
		},
		{
			name: "Paper vs Paper",
			h:    Paper,
			args: args{other: Paper},
			want: Draw,
		},
		{
			name: "Paper vs Scissors",
			h:    Paper,
			args: args{other: Scissors},
			want: Loss,
		},
		{
			name: "Scissors vs Rock",
			h:    Scissors,
			args: args{other: Rock},
			want: Loss,
		},
		{
			name: "Scissors vs Paper",
			h:    Scissors,
			args: args{other: Paper},
			want: Win,
		},
		{
			name: "Scissors vs Scissors",
			h:    Scissors,
			args: args{other: Scissors},
			want: Draw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.h.Outcome(tt.args.other); err != nil {
				t.Errorf("Outcome() error = %v", err)
			} else if got != tt.want {
				t.Errorf("outcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_HandForOutcome(t *testing.T) {
	type args struct {
		outcome  Outcome
		opponent Hand
	}
	tests := []struct {
		name string
		h    Hand
		args args
		want Hand
	}{
		{
			name: "Rock vs Rock",
			h:    Rock,
			args: args{outcome: Draw, opponent: Rock},
			want: Rock,
		},
		{
			name: "Rock vs Paper",
			args: args{outcome: Win, opponent: Rock},
			want: Paper,
		},
		{
			name: "Rock vs Scissors",
			args: args{outcome: Loss, opponent: Rock},
			want: Scissors,
		},
		{
			name: "Paper vs Rock",
			args: args{outcome: Loss, opponent: Paper},
			want: Rock,
		},
		{
			name: "Paper vs Paper",
			args: args{outcome: Draw, opponent: Paper},
			want: Paper,
		},
		{
			name: "Paper vs Scissors",
			args: args{outcome: Win, opponent: Paper},
			want: Scissors,
		},
		{
			name: "Scissors vs Rock",
			args: args{outcome: Win, opponent: Scissors},
			want: Rock,
		},
		{
			name: "Scissors vs Paper",
			args: args{outcome: Loss, opponent: Scissors},
			want: Paper,
		},
		{
			name: "Scissors vs Scissors",
			args: args{outcome: Draw, opponent: Scissors},
			want: Scissors,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := HandForOutcome(tt.args.outcome, tt.args.opponent); err != nil {
				t.Errorf("OpponentHandForOutcome() error = %v", err)
			} else if got != tt.want {
				t.Errorf("OpponentHandForOutcome() = %v, want %v", got, tt.want)
			}

		})
	}
}

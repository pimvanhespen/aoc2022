package main

import "testing"

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
			if got := tt.h.Outcome(tt.args.other); got != tt.want {
				t.Errorf("Outcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_CalcHandForOutcome(t *testing.T) {
	type args struct {
		outcome Outcome
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
			args: args{outcome: Draw},
			want: Rock,
		},
		{
			name: "Rock vs Paper",
			h:    Rock,
			args: args{outcome: Loss},
			want: Paper,
		},
		{
			name: "Rock vs Scissors",
			h:    Rock,
			args: args{outcome: Win},
			want: Scissors,
		},
		{
			name: "Paper vs Rock",
			h:    Paper,
			args: args{outcome: Win},
			want: Rock,
		},
		{
			name: "Paper vs Paper",
			h:    Paper,
			args: args{outcome: Draw},
			want: Paper,
		},
		{
			name: "Paper vs Scissors",
			h:    Paper,
			args: args{outcome: Loss},
			want: Scissors,
		},
		{
			name: "Scissors vs Rock",
			h:    Scissors,
			args: args{outcome: Loss},
			want: Rock,
		},
		{
			name: "Scissors vs Paper",
			h:    Scissors,
			args: args{outcome: Win},
			want: Paper,
		},
		{
			name: "Scissors vs Scissors",
			h:    Scissors,
			args: args{outcome: Draw},
			want: Scissors,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.CalcHandForOutcome(tt.args.outcome); got != tt.want {
				t.Errorf("CalcHandForOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_outcomeFromRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want Outcome
	}{
		{
			name: "Win",
			args: args{r: 'Z'},
			want: Win,
		},
		{
			name: "Loss",
			args: args{r: 'X'},
			want: Loss,
		},
		{
			name: "Draw",
			args: args{r: 'Y'},
			want: Draw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outcomeFromRune(tt.args.r); got != tt.want {
				t.Errorf("outcomeFromRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcRoundScore_Strategy2(t *testing.T) {
	type args struct {
		r Round
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Rock - Draw",
			args: args{r: Round{'A', 'Y'}},
			want: 4,
		},
		{
			name: "Paper - Loss",
			args: args{r: Round{'B', 'X'}},
			want: 1,
		},
		{
			name: "Scissors - Win",
			args: args{r: Round{'C', 'Z'}},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calcRoundScore_Strategy2(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("calcRoundScore_Strategy2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calcRoundScore_Strategy2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

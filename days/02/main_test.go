package main

import (
	"aoc2022/pkg/rps"
	"testing"
)

func Test_outcomeFromRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want rps.Outcome
	}{
		{
			name: "win",
			args: args{r: 'Z'},
			want: rps.Win,
		},
		{
			name: "Loss",
			args: args{r: 'X'},
			want: rps.Loss,
		},
		{
			name: "Draw",
			args: args{r: 'Y'},
			want: rps.Draw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := parseOutcome(tt.args.r); err != nil {
				t.Errorf("parseOutcome() error = %v", err)
			} else if got != tt.want {
				t.Errorf("parseOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolve2(t *testing.T) {
	type args struct {
		r inputRow
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Rock - Draw",
			args: args{r: inputRow{'A', 'Y'}},
			want: 4,
		},
		{
			name: "Paper - Loss",
			args: args{r: inputRow{'B', 'X'}},
			want: 1,
		},
		{
			name: "Scissors - win",
			args: args{r: inputRow{'C', 'Z'}},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := solve([]inputRow{tt.args.r}, transformB, strategyB)

			if (err != nil) != tt.wantErr {
				t.Errorf("strategyB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("strategyB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

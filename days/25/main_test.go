package main

import (
	"fmt"
	"testing"
)

var testInputs = map[int]string{
	1:         "1",
	2:         "2",
	3:         "1=",
	4:         "1-",
	5:         "10",
	6:         "11",
	7:         "12",
	8:         "2=",
	9:         "2-",
	10:        "20",
	15:        "1=0",
	20:        "1-0",
	2022:      "1=11-2",
	12345:     "1-0---0",
	314159265: "1121-1110-1=0",
}

func TestParse(t *testing.T) {

	for expected, input := range testInputs {
		actual, err := ParseSNAFU(input)
		if err != nil {
			t.Errorf("%v", err)
		}
		if actual.decimal != expected {
			t.Errorf("parse(%q) = %d, want %d", input, actual, expected)
		}
	}
}

func TestSnafu_String(t *testing.T) {
	for expected, input := range testInputs {
		snafu := Snafu{expected}
		actual := snafu.String()
		fmt.Println(actual)
		if actual != input {
			t.Errorf("snafu.String() = %q, want %q", actual, input)
		}
	}
}

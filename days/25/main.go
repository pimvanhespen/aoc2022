package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"math"
)

func main() {
	nums, err := aoc.Load(25, parse)
	if err != nil {
		panic(err)
	}

	var sum int
	for _, num := range nums {
		sum += num.decimal
	}

	fmt.Println(sum)
	fmt.Println("Part 1:", Snafu{decimal: sum})
}

func parse(reader io.Reader) ([]Snafu, error) {
	var nums []Snafu
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		num, err := ParseSNAFU(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", text, err)
		}
		nums = append(nums, num)

	}

	return nums, nil
}

type Snafu struct {
	decimal int
}

func ParseSNAFU(input string) (Snafu, error) {
	if len(input) == 0 {
		return Snafu{}, errors.New("input is empty string")
	}

	var sum int

	for i, c := range input {
		var n int
		switch c {
		case '=':
			n = -2
		case '-':
			n = -1
		case '0':
			n = 0
		case '1':
			n = 1
		case '2':
			n = 2
		}

		power := len(input) - 1 - i

		sum += n * int(math.Pow(5, float64(power)))
	}

	return Snafu{sum}, nil
}

func SnafuToDecimal(snafu Snafu) int {
	return snafu.decimal
}

func DecimalToSnafu(decimal int) Snafu {
	return Snafu{decimal}
}

func snafuLimit(power int) int {
	return int(math.Pow(5, float64(power)))
}

var log10of5 = math.Log10(5)

func (s Snafu) String() string {
	var target = float64(s.decimal)

	limit := math.Ceil(math.Log10(target) / log10of5)

	var chars []rune

	var current float64

	for i := limit; i >= 0; i-- {
		if current == target {
			chars = append(chars, toRune(0))
			continue
		}

		var bestMultiplier int

		closestResult := current
		pow := math.Pow(5, i)

		for multiplier := -2; multiplier <= 2; multiplier++ {
			add := float64(multiplier) * pow
			option := current + add
			if math.Abs(target-option) < math.Abs(target-closestResult) {
				bestMultiplier = multiplier
				closestResult = option
			}
		}
		chars = append(chars, toRune(bestMultiplier))
		current = closestResult
	}

	if chars[0] == '0' {
		chars = chars[1:]
	}

	return string(chars)
}

func toRune(n int) rune {
	switch n {
	case -2:
		return '='
	case -1:
		return '-'
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	}
	panic("invalid rune")
}

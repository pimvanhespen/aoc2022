package aoc

import (
	"bufio"
	"fmt"
	"io"
)

type ParseFn[Row any] func(string) (Row, error)

type Parser[Row any] struct {
	SkipEmptyLines bool
	ParseFn        ParseFn[Row]
}

func (p Parser[Row]) Rows(reader io.Reader) ([]Row, error) {
	var rows []Row

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" && p.SkipEmptyLines {
			continue
		}
		row, err := p.ParseFn(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%q': %w", line, err)
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func ParseRows[Row any](reader io.Reader, fn ParseFn[Row]) ([]Row, error) {
	var rows []Row

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			rows = append(rows, *new(Row))
		}
		row, err := fn(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%q': %w", line, err)
		}
		rows = append(rows, row)
	}

	return rows, nil
}

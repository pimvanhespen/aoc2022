package main

import (
	"bytes"
)

// here are several iterations of the same function, with different optimizations
// the "noAllocs" version is the fastest, but it's not very readable

func solve1_lowMem(b []byte) (int, error) {

	left := make(map[byte]bool)
	right := make(map[byte]bool)

	var total int
	var char byte
	var offset int
	var next int

	for {
		next = bytes.Index(b[offset:], []byte{'\n'})

		if next == -1 {
			break
		}

		for i := 0; i < next/2; i++ {
			char = b[offset+i]
			if right[char] {
				total += value(char)
				break
			}
			left[char] = true

			char = b[offset+next-1-i]
			if left[char] {
				total += value(char)
				break
			}
			right[char] = true
		}

		offset += next + 1

		if offset >= len(b) {
			break
		}

		for k := range left {
			left[k] = false
		}

		for k := range right {
			right[k] = false
		}
	}

	return total, nil
}

func solve1_lowMem_IndexByte(b []byte) (int, error) {

	left := make(map[byte]bool)
	right := make(map[byte]bool)

	var total int
	var char byte
	var offset int
	var next int

	const nl = '\n'

	for {
		next = bytes.IndexByte(b[offset:], nl)

		if next == -1 {
			break
		}

		for i := 0; i < next/2; i++ {
			char = b[offset+i]
			if right[char] {
				total += value(char)
				break
			}
			left[char] = true

			char = b[offset+next-1-i]
			if left[char] {
				total += value(char)
				break
			}
			right[char] = true
		}

		offset += next + 1

		if offset >= len(b) {
			break
		}

		for k := range left {
			left[k] = false
		}

		for k := range right {
			right[k] = false
		}
	}

	return total, nil
}

func solve1_lowMem_structMap(b []byte) (int, error) {

	left := make(map[byte]struct{})
	right := make(map[byte]struct{})

	var total int
	var char byte
	var offset int
	var next int

	const nl = '\n'
	var ok bool

	for {
		next = bytes.IndexByte(b[offset:], nl)

		if next == -1 {
			break
		}

		for i := 0; i < next/2; i++ {
			char = b[offset+i]
			if _, ok = right[char]; ok {
				total += value(char)
				break
			} else {
				left[char] = struct{}{}
			}

			char = b[offset+next-1-i]
			if _, ok = left[char]; ok {
				total += value(char)
				break
			} else {
				right[char] = struct{}{}
			}
		}

		offset += next + 1

		if offset >= len(b) {
			break
		}

		for k := range left {
			delete(left, k)
		}

		for k := range right {
			delete(right, k)
		}
	}

	return total, nil
}

func solve1_speed(b []byte) (int, error) {

	var total int
	var char byte
	var offset int
	var next int

	const nl = '\n'
	var ok bool

	for {
		left := make(map[byte]struct{})
		right := make(map[byte]struct{})

		next = bytes.IndexByte(b[offset:], nl)

		if next == -1 {
			break
		}

		for i := 0; i < next/2; i++ {
			char = b[offset+i]
			if _, ok = right[char]; ok {
				total += value(char)
				break
			} else {
				left[char] = struct{}{}
			}

			char = b[offset+next-1-i]
			if _, ok = left[char]; ok {
				total += value(char)
				break
			} else {
				right[char] = struct{}{}
			}
		}

		offset += next + 1

		if offset >= len(b) {
			break
		}
	}

	return total, nil
}

// solve1_noAllocs - Pretty proud of this one...
func solve1_noAllocs(b []byte) (int, error) {

	const newline = '\n'

	var total int

	endOfLine := bytes.IndexByte(b, newline)
	nextNewLine := endOfLine
	halfway := endOfLine / 2

	for index := 0; index < len(b); index++ {

		// check if the current character (b[index]) can be found in the second half of the line )b[halfway:nextLine])
		// if the value of IndexByte is not -1, then the character is found in the second half of the line
		if bytes.IndexByte(b[halfway:nextNewLine], b[index]) != -1 {

			// found char in next half
			total += value(b[index])

			// set index to next char after newlineBytes
			index = nextNewLine + 1

			// stop if the index passed the end of the slice
			if index > len(b)-1 {
				break
			}

			// find the next newlineBytes
			endOfLine = bytes.IndexByte(b[index:], newline)

			nextNewLine = index + endOfLine
			halfway = index + endOfLine/2
		}

	}

	return total, nil
}

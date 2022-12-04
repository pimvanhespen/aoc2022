package main

import (
	"bytes"
)

var newline = []byte{'\n'}

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

func solve1_stdlib(b []byte) (int, error) {

	const newline = '\n'

	var total int

	ni := bytes.IndexByte(b, newline)
	next := ni
	half := ni / 2

	for i := 0; i < len(b); i++ {
		//if i >= half {
		//
		//	d := next - half
		//
		//	panic(fmt.Sprintf("nomatch found in [%d:%d] %s", i-d, next, string(b[i-d:next])))
		//}

		if bytes.IndexByte(b[half:next], b[i]) != -1 {
			// found char in next half
			total += value(b[i])

			// set index to next char after newline
			i = next + 1

			if i > len(b) {
				break
			}

			ni = bytes.IndexByte(b[i:], newline)

			next = i + ni
			half = i + ni/2
		}

	}

	return total, nil
}

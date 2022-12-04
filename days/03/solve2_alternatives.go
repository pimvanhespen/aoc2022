package main

import (
	"bytes"
)

func solve2_speed(rs []byte) (int, error) {

	const newline = '\n'

	var total int

	var past int

	var ok bool

	m := make(map[byte]struct{})
	m2 := make(map[byte]struct{})

	for i := 0; i < len(rs); i++ {

		if rs[i] == newline {
			past++

			if past < 3 {
				continue
			}

			// reset
			m = make(map[byte]struct{})
			m2 = make(map[byte]struct{})
			past = 0

			continue
		}

		switch past {
		case 0:
			if _, ok = m[rs[i]]; !ok {
				m[rs[i]] = struct{}{}
			}

		case 1:
			if _, ok = m[rs[i]]; ok {
				m2[rs[i]] = struct{}{}
				delete(m, rs[i]) // make sure this case it hit only once
			}

		case 2:
			if _, ok = m2[rs[i]]; ok {
				total += value(rs[i])
			}
			past++

		case 3:
			continue
		}

	}

	return total, nil
}

func solve2_mem(rs []byte) (int, error) {

	const newline = '\n'

	var total int

	var past int

	var ok bool

	m := make(map[byte]struct{})
	m2 := make(map[byte]struct{})

	for i := 0; i < len(rs); i++ {

		if rs[i] == newline {
			past++

			if past < 3 {
				continue
			}

			// reset
			for k := range m {
				delete(m, k)
			}
			for k := range m2 {
				delete(m2, k)
			}

			past = 0
			continue
		}

		switch past {
		case 0:
			if _, ok = m[rs[i]]; !ok {
				m[rs[i]] = struct{}{}
			}

		case 1:
			if _, ok = m[rs[i]]; ok {
				m2[rs[i]] = struct{}{}
				delete(m, rs[i]) // make sure this case it hit only once
			}

		case 2:
			if _, ok = m2[rs[i]]; ok {
				total += value(rs[i])
			}
			past++

		case 3:
			continue
		}

	}

	return total, nil
}

func solve2_bytes(b []byte) (int, error) {

	const newline = '\n'
	const notFound = -1

	var total int

	b1 := 0
	e1 := bytes.IndexByte(b, newline)
	b2 := e1 + 1
	e2 := b2 + bytes.IndexByte(b[b2:], newline)
	b3 := e2 + 1
	e3 := b3 + bytes.IndexByte(b[b3:], newline)

	for i := 0; i < len(b); i++ {

		if bytes.IndexByte(b[b2:e2], b[i]) == notFound {
			continue
		}
		if bytes.IndexByte(b[b3:e3], b[i]) == notFound {
			continue
		}

		total += value(b[i])

		b1 = e3 + 1
		e1 = b1 + bytes.IndexByte(b[b1:], newline)
		b2 = e1 + 1
		e2 = b2 + bytes.IndexByte(b[b2:], newline)
		b3 = e2 + 1
		e3 = b3 + bytes.IndexByte(b[b3:], newline)

		i = b1
	}

	return total, nil
}

func nextPtr(b []byte, b1, e1, b2, e2, b3, e3 *int) {
	const newline = '\n'

	*b1 = *e3 + 1
	*e1 = *b1 + bytes.IndexByte(b[*b1:], newline)
	*b2 = *e1 + 1
	*e2 = *b2 + bytes.IndexByte(b[*b2:], newline)
	*b3 = *e2 + 1
	*e3 = *b3 + bytes.IndexByte(b[*b3:], newline)
}

func solve2_bytes_ptr(b []byte) (int, error) {

	const newline = '\n'
	const notFound = -1

	var total int

	var b1, e1, b2, e2, b3, e3 int

	e3 = -1

	nextPtr(b, &b1, &e1, &b2, &e2, &b3, &e3)

	for i := 0; i < len(b); i++ {

		if bytes.IndexByte(b[b2:e2], b[i]) == notFound {
			continue
		}

		if bytes.IndexByte(b[b3:e3], b[i]) == notFound {
			continue
		}

		total += value(b[i])

		nextPtr(b, &b1, &e1, &b2, &e2, &b3, &e3)

		i = b1
	}

	return total, nil
}

func solve2_bytes_ref(b []byte) (int, error) {

	const newline = '\n'
	const notFound = -1

	var total int

	var b1, e1, b2, e2, b3, e3 int

	e3 = -1

	cn := func() {
		b1 = e3 + 1
		e1 = b1 + bytes.IndexByte(b[b1:], newline)
		b2 = e1 + 1
		e2 = b2 + bytes.IndexByte(b[b2:], newline)
		b3 = e2 + 1
		e3 = b3 + bytes.IndexByte(b[b3:], newline)
	}

	cn()

	for i := 0; i < len(b); i++ {

		if bytes.IndexByte(b[b2:e2], b[i]) == notFound {
			continue
		}
		if bytes.IndexByte(b[b3:e3], b[i]) == notFound {
			continue
		}

		total += value(b[i])

		cn()

		i = b1
	}

	return total, nil
}

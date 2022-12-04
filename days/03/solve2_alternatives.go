package main

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

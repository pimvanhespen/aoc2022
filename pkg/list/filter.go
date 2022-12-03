package list

// Filter returns a new slice containing only the elements of the input slice that satisfy the predicate.
func Filter[In any](ins []In, predicate func(In) bool) []In {
	outs := make([]In, 0, len(ins))
	for _, in := range ins {
		if predicate(in) {
			outs = append(outs, in)
		}
	}
	return outs
}

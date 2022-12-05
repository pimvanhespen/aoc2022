package list

type TransformFunc[In, Out any] func(In) Out

func Transform[In, Out any](ins []In, fn TransformFunc[In, Out]) []Out {
	outs := make([]Out, len(ins))
	for i, in := range ins {
		outs[i] = fn(in)
	}
	return outs
}

type TransformErrFunc[In, Out any] func(In) (Out, error)

func TransformErr[In, Out any](ins []In, f TransformErrFunc[In, Out]) ([]Out, error) {
	outs := make([]Out, len(ins))
	for i, in := range ins {
		out, err := f(in)
		if err != nil {
			return nil, err
		}
		outs[i] = out
	}
	return outs, nil
}

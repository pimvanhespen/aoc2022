package list

import "sync"

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

// TransformParallel is like Transform, but runs the transform function in parallel.
// The order of the output slice is guaranteed to be the same as the input slice.
func TransformParallel[In, Out any](ins []In, f TransformFunc[In, Out]) []Out {
	outs := make([]Out, len(ins))
	wg := sync.WaitGroup{}
	wg.Add(len(ins))
	for i, in := range ins {
		go func(i int, in In) {
			defer wg.Done()
			outs[i] = f(in)
		}(i, in)
	}
	wg.Wait()
	return outs
}

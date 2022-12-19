package list

func Reduce[Result any, Item any](items []Item, initial Result, fn func(Result, Item) Result) Result {
	for _, item := range items {
		initial = fn(initial, item)
	}
	return initial
}

func ReduceIndex[Result any, Item any](items []Item, initial Result, fn func(Result, Item, int) Result) Result {
	for i, item := range items {
		initial = fn(initial, item, i)
	}
	return initial
}

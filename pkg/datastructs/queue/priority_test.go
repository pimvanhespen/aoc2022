package queue

import "testing"

type item[T comparable] struct {
	t T
}

func (i item[T]) Equals(t item[T]) bool {
	return i.t == t.t
}

func TestPriority_Pop(t *testing.T) {
	q := NewPriority[item[int]]()
	q.Insert(item[int]{1}, 0)
	q.Insert(item[int]{2}, 1)
	q.Insert(item[int]{3}, 1)

	if i := q.Pop(); i.t != 1 {
		t.Errorf("expect 1, got %d", i.t)
	}

	if i := q.Pop(); i.t != 3 {
		t.Errorf("expect 3, got %d", i.t)
	}

	if i := q.Pop(); i.t != 2 {
		t.Errorf("expect 2, got %d", i.t)
	}
}

func TestPriority_Upsert(t *testing.T) {
	q := NewPriority[item[int]]()
	q.Insert(item[int]{1}, 10)
	q.Insert(item[int]{2}, 10)
	q.Insert(item[int]{3}, 10)
	q.Insert(item[int]{4}, 10)
	q.Insert(item[int]{5}, 10)
	q.Upsert(item[int]{3}, 0)
	q.Upsert(item[int]{2}, 0)

	if i := q.Pop(); i.t != 2 {
		t.Errorf("expected 2, got %d", i.t)
	}

	if i := q.Pop(); i.t != 3 {
		t.Errorf("expected 3, got %d", i.t)
	}

}

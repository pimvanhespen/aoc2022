package stack

import (
	"testing"
)

func BenchmarkStack_Push(b *testing.B) {

	b.Run("Push", func(b *testing.B) {
		s := New[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push(i)
		}
		b.StopTimer()
		reportBad(s)
	})

	b.Run("PushBack", func(b *testing.B) {
		s := New[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.pushFront(i)
		}
		b.StopTimer()
		reportBad(s)
	})
}

func reportBad(s Stack[int]) {
	if s.Len() < 0 {
		panic("bad")
	}
}

func TestStack_Pop(t *testing.T) {
	s := New[int](1, 2, 3)

	if s.Pop() != 3 {
		t.Errorf("expected 3, got %d", s.Pop())
	}
	if s.Pop() != 2 {
		t.Errorf("expected 2, got %d", s.Pop())
	}
	if s.Pop() != 1 {
		t.Errorf("expected 1, got %d", s.Pop())
	}
}

func TestStack_Copy(t *testing.T) {
	s1 := New[int](1, 2, 3)
	s2 := s1.Copy()
	s1.Pop()
	if s1.Len() != 2 {
		t.Errorf("expected 2, got %d", s1.Len())
	}

	if s2.Len() != 3 {
		t.Errorf("expected 3, got %d", s2.Len())
	}
}
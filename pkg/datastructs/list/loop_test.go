package list

import "testing"

func TestLoop_InsertBefore(t *testing.T) {
	loop := NewLoop(1, 2, 3, 4, 5)
	loop.InsertBefore(6)
	if loop.size != 6 {
		t.Errorf("Expected size 6, got %d", loop.size)
	}
	if loop.head.Value != 1 {
		t.Errorf("Expected current value 1, got %d", loop.head.Value)
	}

	if loop.head.prev.Value != 6 {
		t.Errorf("Expected current.prev value 6, got %d", loop.head.prev.Value)
	}
}

func TestLoop_Remove(t *testing.T) {
	loop := NewLoop(1, 2, 3, 4, 5)
	loop.Remove()
	if loop.size != 4 {
		t.Errorf("Expected size 4, got %d", loop.size)
	}
	if loop.head.Value != 2 {
		t.Errorf("Expected current value 2, got %d", loop.head.Value)
	}
}

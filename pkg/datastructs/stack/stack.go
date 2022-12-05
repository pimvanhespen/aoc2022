package stack

import (
	"fmt"
	"strings"
)

// Stack is a stack of items.
// The zero value for Stack is an empty stack ready to use.
type Stack[Item any] []Item

// New returns a new stack with the given items.
// The items are pushed in the order they are given.
func New[Item any](items ...Item) Stack[Item] {
	if len(items) == 0 {
		return Stack[Item]{}
	}

	s := make(Stack[Item], 0, len(items))
	for _, item := range items {
		s.Push(item)
	}
	return s
}

func (s *Stack[Item]) Push(item Item) {
	// by pushing to the back of the slice, we can use the same slice for the stack
	// which is quite a lot faster than pushing to the front of the slice.
	*s = append(*s, item)
}

func (s *Stack[Item]) Pop() Item {
	// Because we push to the back of the slice, we must also pop from the back.
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}

func (s Stack[Item]) Peek() Item {
	return s[len(s)-1]
}

func (s Stack[Item]) String() string {
	var b strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteString(fmt.Sprintf("%v", s[i]))
		if i > 0 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func (s Stack[Item]) Len() int {
	return len(s)
}

func (s Stack[Item]) PeekAt(idx int) Item {
	if idx >= len(s) {
		panic(fmt.Sprintf("index [%d] out of range", idx))
	}
	return s[len(s)-1-idx]
}

// pushFront pushes to the front of the stack. It should be the 'correct' way to
// push to the stack, but it's slower than the other way. So we don't support pushBack in this stack.
// See the BenckmarkStack_Push test for more info.
func (s *Stack[Item]) pushFront(i Item) {
	*s = append([]Item{i}, *s...)
}

func (s Stack[Item]) Copy() Stack[Item] {
	c := make(Stack[Item], len(s))
	copy(c, s)
	return c
}

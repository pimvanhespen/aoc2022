package list

import (
	"fmt"
	"strings"
)

type node[Item any] struct {
	Value Item
	next  *node[Item]
	prev  *node[Item]
}

func newNode[Item any](value Item) *node[Item] {
	return &node[Item]{
		Value: value,
	}
}

type Loop[Item any] struct {
	head *node[Item]
	size int
}

func NewLoop[T any](items ...T) *Loop[T] {
	if len(items) == 0 {
		panic("Cannot create loop with no items")
	}

	head := &node[T]{Value: items[0]}
	tail := head

	for i := 1; i < len(items); i++ {
		n := &node[T]{Value: items[i]}
		tail.next = n
		n.prev = tail
		tail = n
	}

	tail.next = head
	head.prev = tail

	l := &Loop[T]{head: head, size: len(items)}

	return l
}

func (l *Loop[Item]) Size() int {
	return l.size
}

func (l *Loop[Item]) Value() Item {
	return l.head.Value
}

func (l *Loop[Item]) Next() {
	l.head = l.head.next
}

func (l *Loop[Item]) Prev() {
	l.head = l.head.prev
}

func (l *Loop[Item]) Move(n int) {
	if n == 0 {
		return
	}

	for n > 0 {
		l.Next()
		n--
	}

	for n < 0 {
		l.Prev()
		n++
	}
}

func (l *Loop[Item]) String() string {
	if l.size == 0 {
		return "[]"
	}

	if l.size == 1 {
		return fmt.Sprintf("[%v]", l.head.Value)
	}

	var sb strings.Builder
	sb.WriteRune('[')
	sb.WriteString(fmt.Sprintf("%v", l.head.Value))
	n := l.head.next
	for i := 1; i < l.size; i++ {
		sb.WriteString(fmt.Sprintf(", %v", n.Value))
		n = n.next
	}
	sb.WriteByte(']')
	return sb.String()
}

func (l *Loop[Item]) Remove() Item {
	if l.size == 0 {
		panic("Cannot remove from empty loop")
	}

	curr := l.head

	if l.size == 1 {
		l.head = nil
	} else {
		l.Next()
		unlink(curr)
	}

	l.size--
	return curr.Value
}

func (l *Loop[Item]) InsertBefore(item Item) {
	link(l.head.prev, newNode(item), l.head)
	l.size++
}

func unlink[Item any](item *node[Item]) {
	item.prev.next = item.next
	item.next.prev = item.prev
	item.next = nil
	item.prev = nil
}

func link[Item any](prev, item, next *node[Item]) {
	prev.next = item
	item.prev = prev
	item.next = next
	next.prev = item
}

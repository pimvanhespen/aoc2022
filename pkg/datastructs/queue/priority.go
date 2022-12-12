package queue

import (
	"sort"
)

type Equatable[T any] interface {
	Equals(T) bool
}

type node[Item Equatable[Item]] struct {
	Prio  int
	Value Item
}

// Priority[Item any] is a priorityqueue
// Item may be anything (the item type), the priority must be comparable
type Priority[Item Equatable[Item]] struct {
	nodes []node[Item]
}

func NewPriority[Item Equatable[Item]]() *Priority[Item] {
	return &Priority[Item]{}
}

func (p *Priority[Item]) Insert(item Item, priority int) {
	i := sort.Search(len(p.nodes), func(i int) bool {
		return p.nodes[i].Prio >= priority
	})
	n := node[Item]{
		Prio:  priority,
		Value: item,
	}
	p.nodes = append(p.nodes[:i], append([]node[Item]{n}, p.nodes[i:]...)...)
}

func (p *Priority[Item]) Pop() Item {
	top := p.nodes[0]
	p.nodes = p.nodes[1:]
	return top.Value
}

func (p Priority[T]) Len() int {
	return len(p.nodes)
}

func (p *Priority[Item]) Contains(item Item) bool {
	for _, n := range p.nodes {
		if item.Equals(n.Value) {
			return true
		}
	}
	return false
}

func (p *Priority[Item]) Upsert(item Item, prio int) {
	var index int
	var found bool
	for ; index < len(p.nodes); index++ {
		if p.nodes[index].Value.Equals(item) {
			found = true
			break
		}
	}
	if found {
		p.nodes = append(p.nodes[:index], p.nodes[index+1:]...)
	}

	p.Insert(item, prio)
}

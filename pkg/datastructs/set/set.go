package set

import (
	"fmt"
)

type Set[Elem comparable] struct {
	m map[Elem]struct{}
}

func New[Elem comparable](elems ...Elem) *Set[Elem] {
	s := Set[Elem]{
		m: make(map[Elem]struct{}, len(elems)),
	}
	for _, e := range elems {
		s.Add(e)
	}
	return &s
}

func FromSlice[Elem comparable](elems []Elem) *Set[Elem] {
	return New[Elem](elems...)
}

func FromMap[Elem comparable](m map[Elem]struct{}) *Set[Elem] {
	s := Set[Elem]{
		m: make(map[Elem]struct{}, len(m)),
	}
	for e := range m {
		s.Add(e)
	}
	return &s
}

// Add adds an element to the set.
func (s *Set[Elem]) Add(e Elem) {
	if _, ok := s.m[e]; !ok {
		s.m[e] = struct{}{}
	}
}

// Remove an element from the set.
// Note that this does not change the underlying map, so the set will still use the same memory.
func (s *Set[Elem]) Remove(e Elem) {
	delete(s.m, e)
}

// Contains returns true if the set contains the element.
func (s *Set[Elem]) Contains(e Elem) bool {
	_, ok := s.m[e]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set[Elem]) Len() int {
	return len(s.m)
}

// IsEmpty returns true if the set is empty.
func (s *Set[Elem]) IsEmpty() bool {
	return s.Len() == 0
}

// ToSlice returns a slice of the elements in the set.
func (s *Set[Elem]) ToSlice() []Elem {
	elems := make([]Elem, 0, s.Len())
	for e := range s.m {
		elems = append(elems, e)
	}
	return elems
}

// Union returns a new set containing all elements in both sets.
func (s *Set[Elem]) Union(other *Set[Elem]) *Set[Elem] {
	u := New[Elem]()
	for e := range s.m {
		u.Add(e)
	}
	for e := range other.m {
		u.Add(e)
	}
	return u
}

// Intersection returns a new set containing all elements that exist in both sets.
func (s *Set[Elem]) Intersection(other *Set[Elem]) *Set[Elem] {
	i := New[Elem]()
	for e := range s.m {
		if other.Contains(e) {
			i.Add(e)
		}
	}
	return i
}

// Difference returns a new set containing all elements that exist in the first set but not the second.
func (s *Set[Elem]) Difference(other *Set[Elem]) *Set[Elem] {
	d := New[Elem]()
	for e := range s.m {
		if !other.Contains(e) {
			d.Add(e)
		}
	}
	return d
}

// SymmetricDifference returns a new set containing all elements that exist in exactly one of the sets.
func (s *Set[Elem]) SymmetricDifference(other *Set[Elem]) *Set[Elem] {
	return s.Union(other).Difference(s.Intersection(other))
}

// IsSubset returns true if the set is a subset of the other set.
func (s *Set[Elem]) IsSubset(other *Set[Elem]) bool {
	for e := range s.m {
		if !other.Contains(e) {
			return false
		}
	}
	return true
}

// IsSuperset returns true if the set is a superset of the other set.
func (s *Set[Elem]) IsSuperset(other *Set[Elem]) bool {
	return other.IsSubset(s)
}

// Equal returns true if the sets contain the same elements.
func (s *Set[Elem]) Equal(other *Set[Elem]) bool {
	return s.IsSubset(other) && other.IsSubset(s)
}

// Clone returns a copy of the set.
func (s *Set[Elem]) Clone() *Set[Elem] {
	return s.Union(New[Elem]())
}

// Clear removes all elements from the set.
// Note that this does not change the underlying map, so the set will still use the same memory.
func (s *Set[Elem]) Clear() {
	for e := range s.m {
		delete(s.m, e)
	}
}

func (s *Set[Elem]) Pop() Elem {
	for e := range s.m {
		s.Remove(e)
		return e
	}
	panic("set is empty")
}

func (s *Set[Elem]) Filter(fn func(e Elem) (keep bool)) *Set[Elem] {
	res := New[Elem]()
	for e := range s.m {
		if fn(e) {
			res.Add(e)
		}
	}
	return res
}

func (s *Set[Elem]) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

func (s *Set[Elem]) ContainsAny(other *Set[Elem]) bool {
	for e := range s.m {
		if _, ok := other.m[e]; ok {
			return true
		}
	}
	return false
}

package set

import "sync"

type Int64Set struct {
	elements map[int64]void
	mu       sync.Mutex
}

func (s *Int64Set) Add(element int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Contains(element) {
		return false
	}

	s.elements[element] = null
	return true
}

func (s *Int64Set) Delete(element int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Contains(element) {
		delete(s.elements, element)
		return true
	} else {
		return false
	}
}

func (s *Int64Set) Contains(element int64) bool {
	_, ok := s.elements[element]
	return ok
}

func (s *Int64Set) Size() int {
	return len(s.elements)
}

func (s *Int64Set) Iterator() <-chan int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	iterator := make(chan int64, len(s.elements))
	for v, _ := range s.elements {
		iterator <- v
	}
	close(iterator)
	return iterator
}

func (s *Int64Set) ToArr() []int64 {
	arr := make([]int64, 0)
	for e := range s.Iterator() {
		arr = append(arr, e)
	}
	return arr
}

func NewInt64() *Int64Set {
	return &Int64Set{
		elements: make(map[int64]void),
	}
}

func FromInt64Arr(arr []int64) *Int64Set {
	set := NewInt64()
	for _, e := range arr {
		set.Add(e)
	}
	return set
}

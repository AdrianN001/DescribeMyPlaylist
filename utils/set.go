package utils

type Set[T comparable] struct {
	data map[T]struct{}

	len int
}

func SetInit[T comparable]() Set[T] {
	return Set[T]{
		data: map[T]struct{}{},
	}
}

func (set Set[T]) Len() int {
	return set.len
}

func (set *Set[T]) Add(item T) bool {
	if set.Contains(item) {
		return false
	}
	set.data[item] = struct{}{}
	set.len++
	return true
}

func (set *Set[T]) Contains(item T) bool {
	_, ok := set.data[item]
	return ok
}

func (set *Set[T]) Remove(item T) {
	if !set.Contains(item) {
		return
	}
	set.len--
	delete(set.data, item)
}

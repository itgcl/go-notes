package main

type Tree[T any] struct {
	cmp  func(T, T) int
	root *node[T]
}

// A node in a Tree.
type node[T any] struct {
	left, right *node[T]
	val         T
}

func (bt *Tree[T]) find(val T) **node[T] {
	pl := &bt.root
	for *pl != nil {
		switch cmp := bt.cmp(val, (*pl).val); {
		case cmp < 0:
			pl = &(*pl).left
		case cmp > 0:
			pl = &(*pl).right
		default:
			return pl
		}
	}
	return pl
}

func (bt *Tree[T]) Insert(val T) bool {
	pl := bt.find(val)
	if *pl != nil {
		return false
	}
	*pl = &node[T]{val: val}
	return true
}

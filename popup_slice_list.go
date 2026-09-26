package lru

import (
	"iter"
)

type sliceNode[T any] struct {
	prev, next int
	val        T
}

type linkedSlice[T any] struct {
	capacity           int
	length             int
	slice              []sliceNode[T]
	dataHead, dataTail int
	freeHead, freeTail int
}

func newLinkedSlice[T any](capacity int) *linkedSlice[T] {
	return &linkedSlice[T]{
		capacity: capacity,
		slice:    make([]sliceNode[T], 0, capacity),
		dataHead: -1,
		dataTail: -1,
		freeHead: -1,
		freeTail: -1,
	}
}

func (ls *linkedSlice[T]) Len() int {
	return ls.length
}

func (ls *linkedSlice[T]) Front() T {
	return ls.slice[ls.dataHead].val
}

func (ls *linkedSlice[T]) Back() T {
	return ls.slice[ls.dataTail].val
}

func (ls *linkedSlice[T]) Push(val T) {
	if ls.capacity != 0 && ls.length == ls.capacity {
		// уже имеем максимальное число элементов
		ls.Remove(ls.dataTail)
	}
	if ls.freeHead == -1 {
		node := sliceNode[T]{
			prev: -1,
			next: ls.dataHead,
			val:  val,
		}
		idx := len(ls.slice)
		if ls.dataHead != -1 {
			ls.slice[ls.dataHead].prev = idx
		}
		ls.dataHead = idx
		ls.slice = append(ls.slice, node)
	} else {
		idx := ls.freeHead
		ls.slice[ls.dataHead].prev = idx
		ls.slice[idx].next = ls.dataHead
		ls.slice[idx].val = val
		ls.freeHead = ls.slice[idx].next
		if ls.freeHead == -1 {
			ls.slice[ls.freeHead].prev = -1
		}
	}
	ls.length++
}

func (ls *linkedSlice[T]) Remove(idx int) {
	if prev := ls.slice[idx].prev; prev != -1 {
		ls.slice[prev].next = ls.slice[idx].next
	}
	if next := ls.slice[idx].next; next != -1 {
		ls.slice[next].prev = ls.slice[idx].prev
	}
	if ls.dataHead == idx {
		ls.dataHead = ls.slice[idx].next
	}
	if ls.dataTail == idx {
		ls.dataTail = ls.slice[idx].prev
	}
	ls.slice[idx].next = -1
	ls.slice[idx].prev = ls.freeTail
	if ls.freeTail != -1 {
		ls.slice[ls.freeTail].next = idx
	} else {
		ls.freeHead = idx
	}
	ls.freeTail = idx

	ls.length--
}

func (ls *linkedSlice[T]) Values() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		curr := ls.dataHead
		idx := 0
		for curr != -1 {
			if !yield(idx, ls.slice[curr].val) {
				return
			}
			idx++
			curr = ls.slice[curr].next
		}
	}
}

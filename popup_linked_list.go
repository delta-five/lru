package lru

import (
	"iter"
	"sync"
)

type llNode[T any] struct {
	next, prev *llNode[T]
	val        T
}

type linkedList[T any] struct {
	capacity   int
	length     int
	head, tail *llNode[T]
	pool       sync.Pool
}

func newLinkedList[T any](capacity int) *linkedList[T] {
	return &linkedList[T]{
		capacity: capacity,
		pool: sync.Pool{
			New: func() interface{} {
				return new(llNode[T])
			},
		},
	}
}

func (ll *linkedList[T]) Len() int {
	return ll.length
}

func (ll *linkedList[T]) Front() T {
	return ll.head.val
}

func (ll *linkedList[T]) Back() T {
	return ll.tail.val
}

func (ll *linkedList[T]) Push(val T) {
	if ll.capacity != 0 && ll.length == ll.capacity {
		// уже имеем максимальное число элементов
		currTail := ll.tail
		if currTail.prev != nil {
			ll.tail = ll.tail.prev
			ll.tail.next = nil
		}
		ll.pool.Put(currTail)
		ll.length--
	}
	newHead := ll.pool.Get().(llNode[T])
	newHead.val, newHead.next = val, ll.head
	if ll.head != nil {
		ll.head.prev = &newHead
	}
	ll.head = &newHead
	ll.length++

}

func (ll *linkedList[T]) Remove(node *llNode[T]) {
	ll.length--
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if ll.head == node {
		ll.head = node.next
	}
	if ll.tail == node {
		ll.tail = node.prev
	}
}

func (ll *linkedList[T]) Values() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		curr := ll.head
		idx := 0
		for curr != nil {
			if !yield(idx, curr.val) {
				return
			}
			idx++
			curr = curr.next
		}

	}
}

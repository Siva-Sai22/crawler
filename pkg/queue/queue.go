// Package queue provides a simple queue implementation using linked list.
package queue

type node[T any] struct {
	value T
	next  *node[T]
}

type Queue[T any] struct {
	head *node[T]
	tail *node[T]
}

func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (que *Queue[T]) IsEmpty() bool {
	return que.head == nil
}

func (que *Queue[T]) Enqueue(value T) {
	newNode := &node[T]{value, nil}
	if que.IsEmpty() {
		que.head = newNode
		que.tail = newNode
	} else {
		que.tail.next = newNode
		que.tail = newNode
	}
}

func (que *Queue[T]) Dequeue() (T, bool) {
	if que.IsEmpty() {
		var zero T
		return zero, false
	}

	node := que.head
	que.head = node.next
	if que.head == nil {
		que.tail = nil
	}
	return node.value, true
}

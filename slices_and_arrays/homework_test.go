package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularableInts interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[CI CircularableInts] struct {
	values []CI
	front  int
	rear   int
}

func NewCircularQueue[CI CircularableInts](size int) CircularQueue[CI] {
	return CircularQueue[CI]{
		values: make([]CI, size),
		front:  -1,
		rear:   -1,
	}
}

func (q *CircularQueue[CI]) Push(value CI) bool {
	if q.Full() {
		return false
	}

	q.rear = (q.rear + 1) % len(q.values)
	q.values[q.rear] = value

	if q.front == -1 {
		q.front = 0
	}

	return true
}

func (q *CircularQueue[CI]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.values[q.front] = 0

	if q.front == 0 && q.rear == 0 {
		q.front = -1
		q.rear = -1
	} else {
		q.front = (q.front + 1) % len(q.values)
	}
	return true
}

func (q *CircularQueue[CI]) Front() CI {
	if q.Empty() {
		return -1
	}
	res := q.values[q.front]
	return res
}

func (q *CircularQueue[CI]) Back() CI {
	if q.Empty() {
		return -1
	}
	res := q.values[q.rear]
	return res
}

func (q *CircularQueue[CI]) Empty() bool {
	return q.front == -1 && q.rear == -1
}

func (q *CircularQueue[CI]) Full() bool {
	return (q.front == 0 && q.rear == len(q.values)-1) || q.front == q.rear+1
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

package profiles

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewOrderedMap[C comparable](less func(a, b C) bool) OrderedMap[C] {
	return OrderedMap[C]{
		lessFn: less,
	}
}

type OrderedMap[C comparable] struct {
	root *node[C]
	lessFn func(a, b C) bool
}

type node[C comparable] struct {
	key, value  C
	left, right *node[C]
	less        func(a, b C) bool
}

func (m *OrderedMap[C]) Insert(key, value C) {
	m.root = insert(m.root, key, value, m.lessFn)
}

func insert[C comparable](n *node[C], key, value C, less func(a, b C) bool) *node[C] {
	if n == nil {
		return &node[C]{key: key, value: value, left: nil, right: nil, less: less}
	}

	if n.key == key {
		n.value = value
	}

	if !n.less(n.key, key) {
		n.left = insert(n.left, key, value, less)
	}

	if n.less(n.key, key) {
		n.right = insert(n.right, key, value, less)
	}

	return n
}

func (m *OrderedMap[C]) Erase(key C) {
	if m.root == nil {
		return
	}

	erase(m.root, key)
}

func erase[C comparable](n *node[C], key C) *node[C] {
	if n == nil {
		return nil
	}

	if !n.less(n.key, key) && n.key != key {
		n.left = erase(n.left, key)
		return n
	}
	if n.less(n.key, key) {
		n.right = erase(n.right, key)
		return n
	}

	if n.left == nil {
		return n.right
	}
	if n.right == nil {
		return n.left
	}

	min := n.right
	for min != nil && min.left != nil {
		min = min.left
	}
	n.key = min.key
	n.value = min.value
	n.right = erase(n.right, n.key)

	return n
}

func (m *OrderedMap[C]) Contains(key C) bool {
	if m.root == nil {
		return false
	}

	return contains(m.root, key)
}

func contains[C comparable](n *node[C], key C) bool {
	if n == nil {
		return false
	}
	if n.key == key {
		return true
	}

	if !n.less(n.key, key) {
		return contains(n.left, key)
	} else {
		return contains(n.right, key)
	}
}

func (m *OrderedMap[C]) Size() int {
	size := 0
	walkIn(m.root, func(_, _ C) {
		size += 1
	})

	return size
}

func walkIn[C comparable](n *node[C], action func(C, C)) {
	if n == nil {
		return
	}

	walkIn(n.left, action)
	action(n.key, n.value)
	walkIn(n.right, action)
}

func (m *OrderedMap[C]) ForEach(action func(k, v C)) {
	walkIn(m.root, action)
}

func lessFn(v1, v2 int) bool {
	return v1 < v2
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int](lessFn)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
